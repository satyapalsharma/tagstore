package cli

import (
	"fmt"
	"os"

	"tagstore/internal/commands"
	"tagstore/src/lib"

	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import <input.csv>",
	Short: "Import tags from a CSV file",
	Long:  "Imports file-tag associations from a CSV file. The CSV must have a header row with 'file_path' and 'tag_name' columns.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inputPath := args[0]

		rawDB := db.RawDB()
		count, err := commands.ImportFromCSV(rawDB, inputPath)
		if err != nil {
			lib.Get().Error("Failed to import CSV: %v", err)
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Imported %d associations from %s\n", count, inputPath)
		lib.Get().Info("Imported %d associations from %s", count, inputPath)
	},
}
