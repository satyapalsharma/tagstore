package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"tagstore/internal/commands"
	"tagstore/internal/db"
	"tagstore/src/lib"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	// Recover from panics in the main goroutine
	defer func() {
		if r := recover(); r != nil {
		lib.Get().Error("Panic recovered: %v", r)
		os.Exit(1)
	}
	}()
}

func main() {
	// Handle graceful shutdown on OS signals
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		lib.Get().Warn("Received shutdown signal, exiting...")
		os.Exit(0)
	}()

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
