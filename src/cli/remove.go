package cli

import (
	"context"
	"fmt"
	"os"

	"tagstore/src/lib"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <file> <tag>...",
	Short: "Remove one or more tags from a file",
	Long:  "Removes the specified tags from the given file.",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		tags := args[1:]

		ctx := context.Background()
		for _, tag := range tags {
			if err := tagService.RemoveTagFromFile(ctx, filePath, tag); err != nil {
				lib.Get().Error("Failed to remove tag '%s' from file '%s': %v", tag, filePath, err)
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Removed tag '%s' from '%s'\n", tag, filePath)
		}
	},
}
