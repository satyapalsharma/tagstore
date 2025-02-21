package cli

import (
	"context"
	"fmt"
	"os"

	"tagstore/src/lib"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [--and|--or] <tag>...",
	Short: "Search files by tags",
	Long:  "Searches for files matching the specified tags. Use --and for AND logic (all tags) or --or for OR logic (any tag).",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		matchAll, err := cmd.Flags().GetBool("and")
		if err != nil {
			lib.Get().Error("Failed to parse --and flag: %v", err)
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		matchOr, err := cmd.Flags().GetBool("or")
		if err != nil {
			lib.Get().Error("Failed to parse --or flag: %v", err)
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if matchAll && matchOr {
			lib.Get().Error("Cannot use both --and and --or flags")
			fmt.Fprintln(os.Stderr, "Error: cannot use both --and and --or flags")
			os.Exit(1)
		}

		// Default to OR if neither specified
		matchAll = matchAll && !matchOr

		ctx := context.Background()
		paths, err := searchService.SearchFilesByTags(ctx, args, matchAll)
		if err != nil {
			lib.Get().Error("Search failed: %v", err)
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, path := range paths {
			fmt.Println(path)
		}
	},
}

func init() {
	searchCmd.Flags().Bool("and", false, "Match files that have ALL specified tags")
	searchCmd.Flags().Bool("or", false, "Match files that have ANY of the specified tags (default)")
}
