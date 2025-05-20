// sshrestore/restore.go
package sshrestore

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh"
)

type Options struct {
	SSHHost     string
	SSHUser     string
	SSHPassword string

	PGUser     string
	PGPassword string
	PGDatabase string

	BackupFilePath string // local file path to .sql
	RemoteTmpPath  string // remote temp path like "/tmp/restore.sql"
}

func RestorePostgresViaSSH(opts Options) error {
	config := &ssh.ClientConfig{
		User: opts.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(opts.SSHPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", opts.SSHHost, config)
	if err != nil {
		return fmt.Errorf("failed to dial SSH: %w", err)
	}
	defer conn.Close()

	// 1. Upload the backup file to remote server
	if err := uploadFile(conn, opts.BackupFilePath, opts.RemoteTmpPath); err != nil {
		return fmt.Errorf("failed to upload SQL backup: %w", err)
	}

	// 2. Run restore commands
	session, err := conn.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	session.Stderr = os.Stderr
	session.Stdout = os.Stdout

	cmd := fmt.Sprintf(`/bin/bash -c '
		export PGPASSWORD="%s"
		if ! psql -U %s -h localhost -lqt | cut -d \| -f 1 | grep -qw %s; then
			createdb -U %s -h localhost %s
		fi
		psql -U %s -h localhost %s < %s
	'`,
		opts.PGPassword,
		opts.PGUser,
		opts.PGDatabase,
		opts.PGUser,
		opts.PGDatabase,
		opts.PGUser,
		opts.PGDatabase,
		opts.RemoteTmpPath,
	)

	return session.Run(cmd)
}

// Uploads a local file to remote via SCP over SSH
func uploadFile(client *ssh.Client, localPath, remotePath string) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	srcFile, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	session.Stderr = os.Stderr
	session.Stdout = os.Stdout
	stdin, err := session.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		defer stdin.Close()
		io.Copy(stdin, srcFile)
	}()

	return session.Run(fmt.Sprintf("cat > %s", remotePath))
}
