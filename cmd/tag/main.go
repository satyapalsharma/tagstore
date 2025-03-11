package main

import (
	"fmt"
	"os"

	"tagstore/internal/commands"
	"tagstore/internal/db"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: tag <command> [args]")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) != 4 {
			fmt.Fprintln(os.Stderr, "Usage: tag add <file> <tag>")
			os.Exit(1)
		}

		filePath := os.Args[2]
		tagName := os.Args[3]

		dbPath := os.Getenv("TAG_DB_PATH")
		if dbPath == "" {
			dbPath = "tags.db"
		}

		conn, err := db.InitDB(dbPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close()

		if err := commands.AddTag(conn, filePath, tagName); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Added tag \"%s\" to file \"%s\"\n", tagName, filePath)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}
