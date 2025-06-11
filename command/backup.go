package command

import (
	"fmt"
	"log"
	"os"

	"github.com/howeyc/gopass" // External package for password masking
	"github.com/oddishodiqul/backup-db-go.git/sshbackup"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Run restore db",
	Run: func(cmd *cobra.Command, args []string) {
		backup()
		os.Exit(0)
	},
}

func backup() {
	// Prompt for SSH password input
	var sshPassword string
	var dbPassword string
	var err error

	fmt.Print("Enter SSH password: ")
	sshPassword, err = readPassword()

	if err != nil {
		log.Fatalf("Failed to read SSH password: %v", err)
	}

	fmt.Print("Enter DB password: ")
	dbPassword, err = readPassword()

	if err != nil {
		log.Fatalf("Failed to read DB password: %v", err)
	}

	// Proceed with the backup process
	err = sshbackup.DumpPostgresViaSSH(sshbackup.Options{
		SSHHost:     os.Getenv("SSH_HOST"),
		SSHUser:     os.Getenv("SSH_USER"),
		SSHPassword: sshPassword, // Use the masked password here

		PGUser:     os.Getenv("PG_USER"),
		PGPassword: dbPassword,
		PGDatabase: os.Getenv("PG_DATABASE_NAME"),

		OutputFile: os.Getenv("DB_BACKUP"),
	})

	if err != nil {
		log.Fatalf("Backup failed: %v", err)
	}

	log.Println("✅ Backup completed successfully")
}

func readPassword() (string, error) {
	// Use gopass to mask the input for password
	password, err := gopass.GetPasswd()
	if err != nil {
		return "", err
	}

	return string(password), nil
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
