package storage

import (
	"context"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/storage"
)

type GCSStorage struct {
	client *storage.Client
	bucket string
}

func NewGCSStorage(bucket string) (*GCSStorage, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %v", err)
	}

	return &GCSStorage{client: client, bucket: bucket}, nil
}

func (s *GCSStorage) Upload(localPath, remotePath string) (string, error) {
	ctx := context.Background()
	file, err := os.Open(localPath)
	if err != nil {
		return "", fmt.Errorf("unable to open file %v", err)
	}
	defer file.Close()

	wc := s.client.Bucket(s.bucket).Object(remotePath).NewWriter(ctx)
	if _, err = io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}

	return fmt.Sprintf("gs://%s/%s", s.bucket, remotePath), nil
}

func (s *GCSStorage) Download(remotePath, localPath string) (string, error) {
	return "", fmt.Errorf("download not implemented yet")
}

func (s *GCSStorage) ListFiles(prefix string) ([]string, error) {
	return nil, fmt.Errorf("list files not implemented yet")
}
