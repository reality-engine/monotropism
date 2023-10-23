package main

import (
	"log"
	"net/http"
	"os"

	"realityengine.org/m/v2/opendream"
)

func main() {
	store := opendream.NewDataStore()
	handler := opendream.NewEEGHandler(store)

	http.HandleFunc("/api/eeg-text-data", handler.ServeEEGTextData)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server is listening on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
