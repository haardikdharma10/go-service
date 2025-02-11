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
