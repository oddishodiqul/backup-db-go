package sshbackup

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

// Options contains configuration for SSH and pg_dump
type Options struct {
	SSHHost     string // e.g. "remote-server:22"
	SSHUser     string
	SSHPassword string

	PGUser     string
	PGPassword string
	PGDatabase string

	OutputFile string // local output file path
}

// DumpPostgresViaSSH streams pg_dump from remote server to local file
func DumpPostgresViaSSH(opts Options) error {
	config := &ssh.ClientConfig{
		User: opts.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(opts.SSHPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // ⚠️ replace with known_hosts or custom callback in prod
	}

	conn, err := ssh.Dial("tcp", opts.SSHHost, config)
	if err != nil {
		return fmt.Errorf("failed to SSH dial: %w", err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	outFile, err := os.Create(opts.OutputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	session.Stdout = outFile
	session.Stderr = os.Stderr

	cmd := fmt.Sprintf("PGPASSWORD='%s' pg_dump -h localhost -U %s %s", opts.PGPassword, opts.PGUser, opts.PGDatabase)

	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("failed to run remote pg_dump: %w", err)
	}

	return nil
}
