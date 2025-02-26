package main

import (
	"database/sql"
	"fmt"
	"os"

	"tagstore/internal/commands"
	"tagstore/internal/db"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Args) < 3 || os.Args[1] != "list-files" {
		fmt.Fprintf(os.Stderr, "Usage: %s list-files <tag>\n", os.Args[0])
		os.Exit(1)
	}

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s list-files <tag>\n", os.Args[0])
		os.Exit(1)
	}

	tagName := os.Args[2]

	dbPath, err := db.DefaultDBPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	db, err := db.Connect(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	paths, err := commands.ListFilesForTag(db, tagName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	for _, path := range paths {
		fmt.Println(path)
	}
}
