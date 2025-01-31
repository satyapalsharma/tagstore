package cli

import (
	"context"
	"fmt"
	"os"

	"tagstore/src/lib"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search files by tag query",
	Long:  "Searches for files matching a boolean tag query. Use AND (space), OR (|), and NOT (-) operators.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		ctx := context.Background()

		paths, err := searchService.SearchFiles(ctx, query)
		if err != nil {
			lib.Get().Error("Search failed: %v", err)
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if len(paths) == 0 {
			fmt.Println("No files found matching query.")
		} else {
			for _, path := range paths {
				fmt.Println(path)
			}
		}
	},
}
