// Package main contains the entry point for the service.
package main

import (
	"fmt"
	"log"
	"net/http"

	"example/go-service/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Root page of the service")
	})
	mux.HandleFunc("/scan", handlers.ScanHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Server starting on port 8080...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
