// Package models contains the data models for the service.
package models

import "time"

// ScanRequest contains the fields needed to request a scan.
type ScanRequest struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
	Path  string `json:"path,omitempty"`
}

// ScanResult contains the fields returned by a scan.
type ScanResult struct {
	TotalFiles     int       `json:"total_files"`
	ProcessedFiles int       `json:"processed_files"`
	TotalPayloads  int       `json:"total_payloads"`
	Errors         []string  `json:"errors,omitempty"`
	Payloads       []Payload `json:"payloads,omitempty"`
}

// Payload contains the payload info and metadata.
type Payload struct {
	SourceFile    string                 `json:"source_file"`
	ScanTime      time.Time              `json:"scan_time"`
	PayloadResult map[string]interface{} `json:"payload_result"`
}

// QueryResult contains the fields returned by a query.
type QueryResult struct {
	ID             string   `json:"id"`
	Severity       string   `json:"severity"`
	CVSS           float64  `json:"cvss"`
	Status         string   `json:"status"`
	PackageName    string   `json:"package_name"`
	CurrentVersion string   `json:"current_version"`
	FixedVersion   string   `json:"fixed_version"`
	Description    string   `json:"description"`
	PublishedDate  string   `json:"published_date"`
	Link           string   `json:"link"`
	RiskFactors    []string `json:"risk_factors"`
}
