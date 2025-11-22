package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	client *s3.Client
	bucket string
}

func NewS3Storage(bucket string, region string) (*S3Storage, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	client := s3.NewFromConfig(cfg)
	return &S3Storage{client: client, bucket: bucket}, nil
}

func (s *S3Storage) Upload(localPath, remotePath string) (string, error) {
	file, err := os.Open(localPath)
	if err != nil {
		return "", fmt.Errorf("unable to open file %v", err)
	}
	defer file.Close()

	_, err = s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(remotePath),
		Body:   file,
	})
	if err != nil {
		return "", fmt.Errorf("unable to upload file, %v", err)
	}

	return fmt.Sprintf("s3://%s/%s", s.bucket, remotePath), nil
}

func (s *S3Storage) Download(remotePath, localPath string) (string, error) {
	// Basic download implementation
	// In a real scenario, we'd use a downloader manager for large files
	return "", fmt.Errorf("download not implemented yet")
}

func (s *S3Storage) ListFiles(prefix string) ([]string, error) {
	return nil, fmt.Errorf("list files not implemented yet")
}
