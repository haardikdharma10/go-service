// Package svc contains the core logic for the service.
package svc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"example/go-service/db"
	"example/go-service/models"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// ScanRepository scans a GitHub repository for .json files and processes them.
func ScanRepository(req models.ScanRequest) (models.ScanResult, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return models.ScanResult{}, fmt.Errorf("GITHUB_TOKEN environment variable is not set")
	}

	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tokenClient := oauth2.NewClient(context.Background(), tokenService)
	client := github.NewClient(tokenClient)

	query := fmt.Sprintf("extension:json repo:%s/%s", req.Owner, req.Repo)
	if req.Path != "" {
		query += fmt.Sprintf(" path:%s", req.Path)
	}

	opts := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var results *github.CodeSearchResult
	var searchErr error
	err := retryGitHubAPICall(func() error {
		var err error
		results, _, err = client.Search.Code(context.Background(), query, opts)
		if err != nil {
			searchErr = err
			return err
		}
		return nil
	})
	if err != nil {
		return models.ScanResult{}, fmt.Errorf("failed to search code after retries: %v", searchErr)
	}

	response := models.ScanResult{
		TotalFiles: *results.Total,
	}

	jobs := make(chan *github.CodeResult, len(results.CodeResults))
	resultsChan := make(chan models.ScanResult, len(results.CodeResults))

	var wg sync.WaitGroup

	numWorkers := getNumWorkers()
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(client, req.Owner, req.Repo, jobs, resultsChan, &wg)
	}

	for _, result := range results.CodeResults {
		jobs <- &result
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for r := range resultsChan {
		response.ProcessedFiles += r.ProcessedFiles
		response.TotalPayloads += r.TotalPayloads
		response.Payloads = append(response.Payloads, r.Payloads...)
		response.Errors = append(response.Errors, r.Errors...)
	}

	response, err = processAndStoreResults(response)
	if err != nil {
		return models.ScanResult{}, err
	}

	return response, nil
}

func worker(client *github.Client, owner, repo string, jobs <-chan *github.CodeResult, results chan<- models.ScanResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		start := time.Now()
		log.Printf("Goroutine started processing file %s at %s", *job.Path, start.Format(time.RFC3339Nano))
		payloads, err := processFile(client, owner, repo, *job.Path)
		result := models.ScanResult{}
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Error processing %s: %v", *job.Path, err))
		} else {
			result.ProcessedFiles = 1
			result.TotalPayloads = len(payloads)
			result.Payloads = payloads
		}
		end := time.Now()
		log.Printf("Goroutine finished processing file %s at %s (took %v)", *job.Path, end.Format(time.RFC3339Nano), end.Sub(start))
		results <- result
	}
}

func getNumWorkers() int {
	numWorkersStr := os.Getenv("NUM_WORKERS")
	if numWorkersStr == "" {
		return 3
	}
	numWorkers, err := strconv.Atoi(numWorkersStr)
	if err != nil {
		return 3
	}
	return numWorkers
}

func processFile(client *github.Client, owner, repo, path string) ([]models.Payload, error) {
	var content *github.RepositoryContent
	var getContentsErr error
	err := retryGitHubAPICall(func() error {
		var err error
		content, _, _, err = client.Repositories.GetContents(context.Background(), owner, repo, path, nil)
		if err != nil {
			getContentsErr = err
			return err
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get contents after retries: %v", getContentsErr)
	}

	decodedContent, err := content.GetContent()
	if err != nil {
		return nil, err
	}

	var jsonData []map[string]interface{}
	if err := json.Unmarshal([]byte(decodedContent), &jsonData); err != nil {
		return nil, err
	}

	var payloads []models.Payload
	scanTime := time.Now()
	for _, payload := range jsonData {
		payloads = append(payloads, models.Payload{
			SourceFile:    path,
			ScanTime:      scanTime,
			PayloadResult: payload,
		})
	}

	return payloads, nil
}

func processAndStoreResults(response models.ScanResult) (models.ScanResult, error) {
	dbConn := db.GetDB()
	for _, payload := range response.Payloads {
		payloadJSON, err := json.Marshal(payload.PayloadResult)
		if err != nil {
			log.Printf("Error marshaling payload: %v", err)
			continue
		}

		scanResults, ok := payload.PayloadResult["scanResults"].(map[string]interface{})
		if !ok {
			log.Printf("Warning: scanResults not found or not a map in payload")
			continue
		}

		vulnerabilities, ok := scanResults["vulnerabilities"].([]interface{})
		if !ok {
			log.Printf("Warning: vulnerabilities not found or not an array in scanResults")
			continue
		}

		for _, vuln := range vulnerabilities {
			vulnMap, ok := vuln.(map[string]interface{})
			if !ok {
				log.Printf("Warning: vulnerability is not a map")
				continue
			}

			severity, ok := vulnMap["severity"].(string)
			if !ok {
				log.Printf("Warning: Severity not found or not a string in vulnerability")
				severity = "UNKNOWN"
			}

			_, err = dbConn.Exec(
				"INSERT INTO scans (source_file, scan_time, severity, payload) VALUES (?, ?, ?, ?)",
				payload.SourceFile,
				payload.ScanTime,
				severity,
				string(payloadJSON),
			)
			if err != nil {
				log.Printf("Error storing scan result: %v", err)
			}
		}
	}
	return response, nil
}

func retryGitHubAPICall(attempt func() error) error {
	maxAttempts := 3 // 1 initial attempt + 2 retries
	var err error

	for i := 0; i < maxAttempts; i++ {
		err = attempt()
		if err == nil {
			return nil
		}
		log.Printf("GitHub API call failed (attempt %d/%d): %v", i+1, maxAttempts, err)
		if i < maxAttempts-1 {
			time.Sleep(time.Second * 2) // sleep for 2 seconds before retry
		}
	}
	return err
}
