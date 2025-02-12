// Package main contains the entry point for the service.
package main

import (
	"fmt"
	"log"
	"net/http"

	"example/go-service/db"
	"example/go-service/handlers"
)

func main() {
	// Initialize db
	err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.GetDB().Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Root page of the service")
	})
	mux.HandleFunc("/scan", handlers.ScanHandler)
	mux.HandleFunc("/query", handlers.QueryHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Server starting on port 8080...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
