package cli

import (
	"context"
	"fmt"
	"os"

	"tagstore/src/lib"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <file> <tag>...",
	Short: "Add one or more tags to a file",
	Long:  "Adds the specified tags to the given file. Creates the file and tags if they don't exist.",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		tags := args[1:]

		ctx := context.Background()
		for _, tag := range tags {
			if err := tagService.AddTagToFile(ctx, filePath, tag); err != nil {
				lib.Get().Error("Failed to add tag '%s' to file '%s': %v", tag, filePath, err)
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Added tag '%s' to '%s'\n", tag, filePath)
		}
	},
}
