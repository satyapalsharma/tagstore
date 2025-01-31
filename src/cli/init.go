package cli

import (
	"fmt"
	"os"

	"tagstore/src/lib"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the database",
	Long:  "Creates the database file and all required tables.",
	Run: func(cmd *cobra.Command, args []string) {
		// Database is already initialized in Execute()
		fmt.Printf("Database initialized at: %s\n", cfg.DatabasePath)
		lib.Get().Info("Database initialized at: %s", cfg.DatabasePath)
	},
}
