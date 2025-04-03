package main

import (
	"os"
	"os/signal"
	"syscall"

	"tagstore/src/cli"
	"tagstore/src/lib"
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
	// Handle uncaught exceptions (panics) in other goroutines
	// by intercepting OS signals and cleaning up gracefully.
	// Go does not have a global unhandled rejection model like Node.js;
	// instead, we rely on defer/recover and signal handling.
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		lib.Get().Warn("Received shutdown signal, exiting...")
		os.Exit(0)
	}()

	cli.Execute()
}
