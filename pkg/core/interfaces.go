package core

// Config holds the configuration for a backup operation
type Config map[string]interface{}

// Database interface for backup and restore operations
type Database interface {
	Backup(config Config, outputPath string) (string, error)
	Restore(config Config, backupPath string) error
	TestConnection(config Config) error
}

// Storage interface for uploading and downloading backups
type Storage interface {
	Upload(localPath, remotePath string) (string, error)
	Download(remotePath, localPath string) (string, error)
	ListFiles(prefix string) ([]string, error)
}
