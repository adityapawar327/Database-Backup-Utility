package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"db-backup-tool/pkg/core"
	"db-backup-tool/pkg/databases"
	"db-backup-tool/pkg/storage"
	"db-backup-tool/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "backup-tool",
	Short: "A CLI tool for database backups",
	Long:  `A robust CLI utility to backup and restore various databases with support for local and cloud storage.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	utils.InitLogger()
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./db_backup_config.yaml)")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(restoreCmd)
	rootCmd.AddCommand(listCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("db_backup_config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		// silently load config
	}
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration",
	Run: func(cmd *cobra.Command, args []string) {
		// In a real app, this would write a default config file
		fmt.Println("Configuration initialized.")
	},
}

var backupCmd = &cobra.Command{
	Use:   "backup [db_name]",
	Short: "Backup a database",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dbName := args[0]
		utils.LogInfo(fmt.Sprintf("Starting backup for %s", dbName))

		if !viper.IsSet(fmt.Sprintf("databases.%s", dbName)) {
			utils.LogError(fmt.Sprintf("Database config '%s' not found", dbName))
			return
		}
		dbConfig := viper.GetStringMap(fmt.Sprintf("databases.%s", dbName))

		dbAdapter, err := getDatabaseAdapter(dbConfig["type"].(string))
		if err != nil {
			utils.LogError(err.Error())
			return
		}

		storageAdapter, err := getStorageAdapter()
		if err != nil {
			utils.LogError(err.Error())
			return
		}

		// Generate temp file name
		ext := "sql"
		if dbConfig["type"].(string) == "sqlite" {
			ext = "db"
		}
		tempFile := fmt.Sprintf("temp_%s_%s.%s", dbName, time.Now().Format("20060102_150405"), ext)

		slackWebhook := viper.GetString("notifications.slack_webhook")

		// Perform Backup
		backupPath, err := dbAdapter.Backup(dbConfig, tempFile)
		if err != nil {
			handleError(fmt.Sprintf("Backup failed for %s: %v", dbName, err), slackWebhook)
			return
		}
		utils.LogInfo(fmt.Sprintf("Database backed up locally to: %s", backupPath))

		// Compress
		compressedPath, err := utils.CompressFile(backupPath)
		if err != nil {
			handleError(fmt.Sprintf("Compression failed: %v", err), slackWebhook)
			os.Remove(backupPath)
			return
		}
		os.Remove(backupPath)
		backupPath = compressedPath
		utils.LogInfo(fmt.Sprintf("Backup compressed to: %s", backupPath))

		// Upload
		storagePath := viper.GetString("storage.path")
		remotePath := filepath.Join(storagePath, filepath.Base(backupPath))
		uploadedPath, err := storageAdapter.Upload(backupPath, remotePath)
		if err != nil {
			handleError(fmt.Sprintf("Upload failed for %s: %v", dbName, err), slackWebhook)
			os.Remove(backupPath)
			return
		}

		successMsg := fmt.Sprintf("Backup successful for %s. Uploaded to: %s", dbName, uploadedPath)
		utils.LogInfo(successMsg)
		utils.SendSlackNotification(slackWebhook, successMsg)

		os.Remove(backupPath)
	},
}

var restoreCmd = &cobra.Command{
	Use:   "restore [backup_file] [db_name]",
	Short: "Restore a database from a backup",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		backupFile := args[0]
		dbName := args[1]
		fmt.Printf("Restoring %s to %s...\n", backupFile, dbName)

		if !viper.IsSet(fmt.Sprintf("databases.%s", dbName)) {
			fmt.Printf("Database config '%s' not found\n", dbName)
			return
		}
		dbConfig := viper.GetStringMap(fmt.Sprintf("databases.%s", dbName))

		dbAdapter, err := getDatabaseAdapter(dbConfig["type"].(string))
		if err != nil {
			fmt.Println(err)
			return
		}

		storageAdapter, err := getStorageAdapter()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Download
		localBackupPath := filepath.Join(os.TempDir(), filepath.Base(backupFile))
		downloadedPath, err := storageAdapter.Download(backupFile, localBackupPath)
		if err != nil {
			fmt.Printf("Download failed: %v\n", err)
			return
		}
		fmt.Printf("Backup downloaded to: %s\n", downloadedPath)

		// Decompress
		restorePath := downloadedPath
		if filepath.Ext(downloadedPath) == ".gz" {
			fmt.Println("Decompressing backup...")
			decompressedPath, err := utils.DecompressFile(downloadedPath)
			if err != nil {
				fmt.Printf("Decompression failed: %v\n", err)
				return
			}
			restorePath = decompressedPath
			fmt.Printf("Decompressed to: %s\n", restorePath)
		}

		// Restore
		if err := dbAdapter.Restore(dbConfig, restorePath); err != nil {
			fmt.Printf("Restore failed: %v\n", err)
			return
		}
		fmt.Println("Database restored successfully!")

		// Cleanup
		if viper.GetString("storage.type") != "local" {
			os.Remove(downloadedPath)
			if restorePath != downloadedPath {
				os.Remove(restorePath)
			}
		}
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List backups",
	Run: func(cmd *cobra.Command, args []string) {
		storageAdapter, err := getStorageAdapter()
		if err != nil {
			fmt.Println(err)
			return
		}

		path := viper.GetString("storage.path")
		files, err := storageAdapter.ListFiles(path)
		if err != nil {
			fmt.Printf("Failed to list files: %v\n", err)
			return
		}

		for _, f := range files {
			fmt.Println(f)
		}
	},
}

func main() {
	Execute()
}

// Helper functions

func getDatabaseAdapter(dbType string) (core.Database, error) {
	switch dbType {
	case "mysql":
		return &databases.MySQLDatabase{}, nil
	case "sqlite":
		return &databases.SQLiteDatabase{}, nil
	case "postgres":
		return &databases.PostgresDatabase{}, nil
	case "mongo":
		return &databases.MongoDatabase{}, nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}

func getStorageAdapter() (core.Storage, error) {
	storageType := viper.GetString("storage.type")
	switch storageType {
	case "local":
		return &storage.LocalStorage{}, nil
	case "s3":
		bucket := viper.GetString("storage.bucket")
		region := viper.GetString("storage.region")
		return storage.NewS3Storage(bucket, region)
	case "gcs":
		bucket := viper.GetString("storage.bucket")
		return storage.NewGCSStorage(bucket)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", storageType)
	}
}

func handleError(msg string, webhook string) {
	utils.LogError(msg)
	utils.SendSlackNotification(webhook, msg)
}
