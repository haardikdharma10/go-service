// Package svc contains the core logic for the service.
package svc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"example/go-service/db"
	"example/go-service/models"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// ScanRepository scans a GitHub repository for .json files and processes them.
func ScanRepository(req models.ScanRequest) (models.ScanResult, error) {
	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "<token>"},
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

	results, _, err := client.Search.Code(context.Background(), query, opts)
	if err != nil {
		return models.ScanResult{}, err
	}

	response := models.ScanResult{
		TotalFiles: *results.Total,
	}

	for _, result := range results.CodeResults {
		payloads, err := processFile(client, req.Owner, req.Repo, *result.Path)
		if err != nil {
			response.Errors = append(response.Errors, fmt.Sprintf("Error processing %s: %v", *result.Path, err))
			continue
		}
		response.ProcessedFiles++
		response.TotalPayloads += len(payloads)
		response.Payloads = append(response.Payloads, payloads...)
	}

	response, err = processAndStoreResults(response)
	if err != nil {
		return models.ScanResult{}, err
	}

	return response, nil
}

func processFile(client *github.Client, owner, repo, path string) ([]models.Payload, error) {
	content, _, _, err := client.Repositories.GetContents(context.Background(), owner, repo, path, nil)
	if err != nil {
		return nil, err
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
