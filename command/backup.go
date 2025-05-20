package command

import (
	"log"
	"os"

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
	err := sshbackup.DumpPostgresViaSSH(sshbackup.Options{
		SSHHost:     os.Getenv("SSH_HOST"),
		SSHUser:     os.Getenv("SSH_USER"),
		SSHPassword: os.Getenv("SSH_PASSWORD"),

		PGUser:     os.Getenv("PG_USER"),
		PGPassword: os.Getenv("PG_PASSWORD"),
		PGDatabase: os.Getenv("PG_DATABASE_NAME"),

		OutputFile: os.Getenv("DB_BACKUP"),
	})

	if err != nil {
		log.Fatalf("Backup failed: %v", err)
	}

	log.Println("✅ Backup completed successfully")
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
