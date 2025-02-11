// Package svc contains the core scanning logic for the service.
package svc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"example/go-service/models"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// ScanRepository scans a GitHub repository for .json files and processes them.
func ScanRepository(req models.ScanRequest) (models.ScanResult, error) {
	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "<github-api-token-here>"},
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
