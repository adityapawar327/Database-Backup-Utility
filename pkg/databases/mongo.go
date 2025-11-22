package databases

import (
	"db-backup-tool/pkg/core"
	"fmt"
	"os/exec"
)

type MongoDatabase struct{}

func (db *MongoDatabase) Backup(config core.Config, outputPath string) (string, error) {
	// mongodump --uri="mongodb://[user]:[password]@[host]:[port]/[database]" --archive=[outputPath]

	user := config["user"].(string)
	password := config["password"].(string)
	host := config["host"].(string)
	port := config["port"].(int)
	database := config["database"].(string)

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", user, password, host, port, database)

	cmd := exec.Command("mongodump",
		fmt.Sprintf("--uri=%s", uri),
		fmt.Sprintf("--archive=%s", outputPath),
	)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("mongodump failed: %v", err)
	}

	return outputPath, nil
}

func (db *MongoDatabase) Restore(config core.Config, backupPath string) error {
	// mongorestore --uri="mongodb://[user]:[password]@[host]:[port]/[database]" --archive=[backupPath]

	user := config["user"].(string)
	password := config["password"].(string)
	host := config["host"].(string)
	port := config["port"].(int)
	database := config["database"].(string)

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", user, password, host, port, database)

	cmd := exec.Command("mongorestore",
		fmt.Sprintf("--uri=%s", uri),
		fmt.Sprintf("--archive=%s", backupPath),
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mongorestore failed: %v", err)
	}

	return nil
}

func (db *MongoDatabase) TestConnection(config core.Config) error {
	return nil // Dummy implementation
}
