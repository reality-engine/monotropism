package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"realityengine.org/m/v2/opendream"

	"cloud.google.com/go/storage"
)

const (
	bucketName = "eeg-raw"                                          // Updated bucket name
	objectName = "eeg-image/EEG_motor_imagery_dataset_edf_only.csv" // Updated object name
)

func main() {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create new storage client: %v", err)
	}
	defer client.Close()

	provider := opendream.NewStorageProvider(client, bucketName, objectName)
	handler := opendream.NewHandler(provider)

	http.HandleFunc("/", handler.ServeCSV)
	// App Engine provides the port to bind to as an environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	log.Printf("Server is listening on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
