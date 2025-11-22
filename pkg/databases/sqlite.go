package databases

import (
	"db-backup-tool/pkg/core"
	"fmt"
	"io"
	"os"
)

type SQLiteDatabase struct{}

func (db *SQLiteDatabase) Backup(config core.Config, outputPath string) (string, error) {
	// SQLite backup is just copying the file
	dbPath := config["path"].(string)

	src, err := os.Open(dbPath)
	if err != nil {
		return "", fmt.Errorf("failed to open sqlite db: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create backup file: %v", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy sqlite db: %v", err)
	}

	return outputPath, nil
}

func (db *SQLiteDatabase) Restore(config core.Config, backupPath string) error {
	// Restore is just copying back
	dbPath := config["path"].(string)

	src, err := os.Open(backupPath)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(dbPath)
	if err != nil {
		return fmt.Errorf("failed to open destination db: %v", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to restore sqlite db: %v", err)
	}

	return nil
}

func (db *SQLiteDatabase) TestConnection(config core.Config) error {
	dbPath := config["path"].(string)
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return fmt.Errorf("sqlite database file does not exist: %s", dbPath)
	}
	return nil
}
