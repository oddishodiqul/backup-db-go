package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/oddishodiqul/backup-db-go.git/command"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Error loading .env file")
	}

	command.Execute()
}
