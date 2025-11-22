package databases

import (
	"db-backup-tool/pkg/core"
	"fmt"
	"os"
	"os/exec"
)

type MySQLDatabase struct{}

func (db *MySQLDatabase) Backup(config core.Config, outputPath string) (string, error) {
	// Construct mysqldump command
	// mysqldump -u [user] -p[password] -h [host] -P [port] [database] > [outputPath]

	user := config["user"].(string)
	password := config["password"].(string)
	host := config["host"].(string)
	port := config["port"].(int)
	database := config["database"].(string)

	cmd := exec.Command("mysqldump",
		fmt.Sprintf("-u%s", user),
		fmt.Sprintf("-p%s", password),
		fmt.Sprintf("-h%s", host),
		fmt.Sprintf("-P%d", port),
		database,
	)

	// Create output file
	outfile, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer outfile.Close()

	cmd.Stdout = outfile

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("mysqldump failed: %v", err)
	}

	return outputPath, nil
}

func (db *MySQLDatabase) Restore(config core.Config, backupPath string) error {
	// mysql -u [user] -p[password] -h [host] -P [port] [database] < [backupPath]

	user := config["user"].(string)
	password := config["password"].(string)
	host := config["host"].(string)
	port := config["port"].(int)
	database := config["database"].(string)

	cmd := exec.Command("mysql",
		fmt.Sprintf("-u%s", user),
		fmt.Sprintf("-p%s", password),
		fmt.Sprintf("-h%s", host),
		fmt.Sprintf("-P%d", port),
		database,
	)

	infile, err := os.Open(backupPath)
	if err != nil {
		return err
	}
	defer infile.Close()

	cmd.Stdin = infile

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mysql restore failed: %v", err)
	}

	return nil
}

func (db *MySQLDatabase) TestConnection(config core.Config) error {
	// mysqladmin -u [user] -p[password] -h [host] -P [port] ping
	return nil // Dummy implementation
}
