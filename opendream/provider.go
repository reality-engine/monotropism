package opendream

import (
	"context"
	"io"
	"net/http"

	"cloud.google.com/go/storage"
)

type Provider interface {
	ServeCSV(w http.ResponseWriter, r *http.Request) error
}

type storageProvider struct {
	client     *storage.Client
	bucketName string
	objectName string
}

func NewStorageProvider(client *storage.Client, bucketName, objectName string) Provider {
	return &storageProvider{client: client, bucketName: bucketName, objectName: objectName}
}

func (s *storageProvider) ServeCSV(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()
	objHandle := s.client.Bucket(s.bucketName).Object(s.objectName)
	reader, err := objHandle.NewReader(ctx)
	if err != nil {
		return err
	}
	defer reader.Close()

	w.Header().Set("Content-Type", "text/csv")
	_, err = io.Copy(w, reader)
	return err
}
