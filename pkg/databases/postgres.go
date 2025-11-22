package databases

import (
	"db-backup-tool/pkg/core"
	"fmt"
	"os"
	"os/exec"
)

type PostgresDatabase struct{}

func (db *PostgresDatabase) Backup(config core.Config, outputPath string) (string, error) {
	// pg_dump -U [user] -h [host] -p [port] [database] > [outputPath]
	// Password is usually supplied via PGPASSWORD env var

	user := config["user"].(string)
	password := config["password"].(string)
	host := config["host"].(string)
	port := config["port"].(int)
	database := config["database"].(string)

	cmd := exec.Command("pg_dump",
		fmt.Sprintf("-U%s", user),
		fmt.Sprintf("-h%s", host),
		fmt.Sprintf("-p%d", port),
		database,
	)

	// Set PGPASSWORD environment variable for this command
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", password))

	outfile, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer outfile.Close()

	cmd.Stdout = outfile

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("pg_dump failed: %v", err)
	}

	return outputPath, nil
}

func (db *PostgresDatabase) Restore(config core.Config, backupPath string) error {
	// psql -U [user] -h [host] -p [port] -d [database] -f [backupPath]

	user := config["user"].(string)
	password := config["password"].(string)
	host := config["host"].(string)
	port := config["port"].(int)
	database := config["database"].(string)

	cmd := exec.Command("psql",
		fmt.Sprintf("-U%s", user),
		fmt.Sprintf("-h%s", host),
		fmt.Sprintf("-p%d", port),
		"-d", database,
		"-f", backupPath,
	)

	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", password))

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("psql restore failed: %v", err)
	}

	return nil
}

func (db *PostgresDatabase) TestConnection(config core.Config) error {
	return nil // Dummy implementation
}
