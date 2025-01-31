package cli

import (
	"context"
	"fmt"
	"os"

	"tagstore/src/lib"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list <file|tag> <name>",
	Short: "List tags for a file or files for a tag",
	Long:  "Lists all tags associated with a file, or all files associated with a tag.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		entityType := args[0]
		name := args[1]

		switch entityType {
		case "file":
			tags, err := tagService.ListTagsForFile(ctx, name)
			if err != nil {
				lib.Get().Error("Failed to list tags for file: %v", err)
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			if len(tags) == 0 {
				fmt.Printf("No tags found for '%s'\n", name)
			} else {
				for _, tag := range tags {
					fmt.Println(tag.Name)
				}
			}
		case "tag":
			files, err := tagService.ListFilesForTag(ctx, name)
			if err != nil {
				lib.Get().Error("Failed to list files for tag: %v", err)
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			if len(files) == 0 {
				fmt.Printf("No files found for tag '%s'\n", name)
			} else {
				for _, file := range files {
					fmt.Println(file.Path)
				}
			}
		default:
			fmt.Fprintf(os.Stderr, "Unknown entity type: %s (use 'file' or 'tag')\n", entityType)
			os.Exit(1)
		}
	},
}
