package main

import (
	"fmt"
	"os"

	"tagstore/internal/commands"
	"tagstore/internal/db"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "remove":
		handleRemove(args)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: tag <command> [args...]\n")
	fmt.Fprintf(os.Stderr, "Commands: remove\n")
}

func handleRemove(args []string) {
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: tag remove <file> <tag>\n")
		os.Exit(1)
	}

	filePath := args[0]
	tagName := args[1]

	dbPath, err := db.DefaultDBPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	sqlDB, err := db.Connect(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	if err := commands.RemoveTag(sqlDB, filePath, tagName); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Removed tag \"%s\" from file \"%s\"\n", tagName, filePath)
}
