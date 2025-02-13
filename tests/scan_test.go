package tests

import (
	"database/sql"
	"testing"

	"example/go-service/db"
	"example/go-service/models"
	"example/go-service/svc"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
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

func TestScanRepository(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	db.SetDB(testDB)

	// Test ScanRepository
	req := models.ScanRequest{
		Owner: "velancio",
		Repo:  "vulnerability_scans",
	}

	result, err := svc.ScanRepository(req)
	if err != nil {
		t.Errorf("ScanRepository returned an error: %v", err)
	}

	// Add assertions to check the result
	if result.TotalFiles != 9 {
		t.Errorf("Expected TotalFiles to be 9, got %d", result.TotalFiles)
	}

	if result.ProcessedFiles != 9 {
		t.Errorf("Expected ProcessedFiles to be 9, got %d", result.ProcessedFiles)
	}
}
