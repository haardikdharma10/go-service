package tests

import (
	"database/sql"
	"encoding/json"
	"testing"

	"example/go-service/db"
	"example/go-service/svc"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDBForQuery(t *testing.T) *sql.DB {
	testDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening test database: %v", err)
	}

	_, err = testDB.Exec(`
		CREATE TABLE scans (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			source_file TEXT,
			scan_time DATETIME,
			severity TEXT,
			payload TEXT
		)
	`)
	if err != nil {
		t.Fatalf("Error creating test table: %v", err)
	}

	return testDB
}

func insertTestData(t *testing.T, testDB *sql.DB) {
	payload := map[string]interface{}{
		"scanResults": map[string]interface{}{
			"vulnerabilities": []map[string]interface{}{
				{
					"id":              "CVE-2021-1234",
					"severity":        "HIGH",
					"cvss":            8.5,
					"status":          "fixed",
					"package_name":    "openssl",
					"current_version": "1.1.1t-r0",
					"fixed_version":   "1.1.1u-r0",
					"description":     "Buffer overflow vulnerability in OpenSSL",
					"published_date":  "2021-01-15T00:00:00Z",
					"link":            "https://nvd.nist.gov/vuln/detail/CVE-2021-1234",
					"risk_factors":    []string{"Remote Code Execution", "High CVSS Score"},
				},
			},
		},
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Error marshaling test payload: %v", err)
	}

	_, err = testDB.Exec("INSERT INTO scans (source_file, scan_time, severity, payload) VALUES (?, ?, ?, ?)",
		"testfile.json", "2023-01-01T00:00:00Z", "HIGH", string(payloadJSON))
	if err != nil {
		t.Fatalf("Error inserting test data: %v", err)
	}
}

func TestQueryVulnerabilities(t *testing.T) {
	testDB := setupTestDBForQuery(t)
	defer testDB.Close()

	db.SetDB(testDB)
	insertTestData(t, testDB)

	tests := []struct {
		name          string
		severity      string
		expectedCount int
		expectedError bool
	}{
		{"Query HIGH severity", "HIGH", 1, false},
		{"Query LOW severity", "LOW", 0, false},
		{"Query invalid severity", "INVALID", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := svc.QueryVulnerabilities(tt.severity)

			if (err != nil) != tt.expectedError {
				t.Errorf("QueryVulnerabilities() error = %v, expectedError %v", err, tt.expectedError)
				return
			}

			if len(results) != tt.expectedCount {
				t.Errorf("QueryVulnerabilities() returned %d results, expected %d", len(results), tt.expectedCount)
			}

			if len(results) > 0 {
				result := results[0]
				if result.Severity != tt.severity {
					t.Errorf("QueryVulnerabilities() returned severity %s, expected %s", result.Severity, tt.severity)
				}
				if result.ID != "CVE-2021-1234" {
					t.Errorf("QueryVulnerabilities() returned ID %s, expected CVE-2021-1234", result.ID)
				}
				if result.CVSS != 8.5 {
					t.Errorf("QueryVulnerabilities() returned CVSS %f, expected 8.5", result.CVSS)
				}
				if len(result.RiskFactors) != 2 {
					t.Errorf("QueryVulnerabilities() returned %d risk factors, expected 2", len(result.RiskFactors))
				}
			}
		})
	}
}
