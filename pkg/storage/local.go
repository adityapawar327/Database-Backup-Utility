package storage

import (
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct{}

func (s *LocalStorage) Upload(localPath, remotePath string) (string, error) {
	// For local storage, "upload" is just copying to the destination directory
	// remotePath here is treated as the destination directory or full path

	destPath := remotePath
	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return "", err
	}

	src, err := os.Open(localPath)
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return destPath, nil
}

func (s *LocalStorage) Download(remotePath, localPath string) (string, error) {
	// For local storage, "download" is copying from the source
	src, err := os.Open(remotePath)
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(localPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return localPath, nil
}

func (s *LocalStorage) ListFiles(prefix string) ([]string, error) {
	var files []string
	err := filepath.Walk(prefix, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
