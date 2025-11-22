package utils

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

func CompressFile(sourcePath string) (string, error) {
	destPath := sourcePath + ".gz"

	src, err := os.Open(sourcePath)
	if err != nil {
		return "", fmt.Errorf("failed to open source file: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to create dest file: %v", err)
	}
	defer dst.Close()

	zw := gzip.NewWriter(dst)
	defer zw.Close()

	if _, err := io.Copy(zw, src); err != nil {
		return "", fmt.Errorf("failed to compress file: %v", err)
	}

	return destPath, nil
}

func DecompressFile(sourcePath string) (string, error) {
	// Assumes sourcePath ends with .gz
	if len(sourcePath) < 3 || sourcePath[len(sourcePath)-3:] != ".gz" {
		return "", fmt.Errorf("invalid file extension for decompression")
	}
	destPath := sourcePath[:len(sourcePath)-3]

	src, err := os.Open(sourcePath)
	if err != nil {
		return "", fmt.Errorf("failed to open source file: %v", err)
	}
	defer src.Close()

	zr, err := gzip.NewReader(src)
	if err != nil {
		return "", fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer zr.Close()

	dst, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to create dest file: %v", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, zr); err != nil {
		return "", fmt.Errorf("failed to decompress file: %v", err)
	}

	return destPath, nil
}
