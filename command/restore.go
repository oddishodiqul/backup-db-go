package command

import (
	"log"
	"os"

	"github.com/oddishodiqul/backup-db-go.git/sshrestore"
	"github.com/spf13/cobra"
)

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Run restore db",
	Run: func(cmd *cobra.Command, args []string) {
		restore()
		os.Exit(0)
	},
}

func restore() {
	err := sshrestore.RestorePostgresViaSSH(sshrestore.Options{
		SSHHost:     os.Getenv("SSH_HOST"),
		SSHUser:     os.Getenv("SSH_USER"),
		SSHPassword: os.Getenv("SSH_PASSWORD"),

		PGUser:     os.Getenv("PG_USER"),
		PGPassword: os.Getenv("PG_PASSWORD"),
		PGDatabase: os.Getenv("PG_DATABASE_TARGET"),

		BackupFilePath: os.Getenv("DB_RESTORE"),
		RemoteTmpPath:  os.Getenv("REMOTE_TMP_PATH"),
	})

	if err != nil {
		log.Fatalf("❌ Restore failed: %v", err)
	}

	log.Println("✅ Restore completed successfully")
}

func init() {
	rootCmd.AddCommand(restoreCmd)
}
