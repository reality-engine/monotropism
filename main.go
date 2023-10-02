package main

import (
	"context"
	"log"
	"net/http"

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

	provider := NewStorageProvider(client, bucketName, objectName)
	handler := NewHandler(provider)

	http.HandleFunc("/", handler.ServeCSV)
	log.Println("Server is listening on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
