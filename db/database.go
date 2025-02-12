// Package db contains the database logic for the service.
package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // just the driver
)

var dbConn *sql.DB

func InitDB() error {
	var err error
	dbConn, err = sql.Open("sqlite3", "./scans.db")
	if err != nil {
		log.Fatal(err)
	}

	scansTable := `
    CREATE TABLE IF NOT EXISTS scans (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        source_file TEXT,
        scan_time DATETIME,
        severity TEXT,
        payload TEXT
    );`

	_, err = dbConn.Exec(scansTable)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func GetDB() *sql.DB {
	return dbConn
}
