package cli

import (
	"context"
	"encoding/csv"
	"database/sql"
	"fmt"
	"os"

	"tagstore/src/lib"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export <output.csv>",
	Short: "Export tag associations to a CSV file",
	Long:  "Exports all file-tag associations to a CSV file with columns file_path and tag_name.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		outputPath := args[0]
		ctx := context.Background()

		if err := ExportToCSV(db.DB(), outputPath); err != nil {
			lib.Get().Error("Failed to export tags: %v", err)
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Exported associations to %s\n", outputPath)
	},
}

// ExportToCSV exports all file-tag associations to a CSV file
func ExportToCSV(db *sql.DB, outputPath string) error {
	if outputPath == "" {
		return fmt.Errorf("output path cannot be empty")
	}

	query := `SELECT f.path, t.name FROM file_tags ft 
		JOIN files f ON ft.file_id = f.id 
		JOIN tags t ON ft.tag_id = t.id 
		ORDER BY f.path, t.name`

	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query tags: %w", err)
	}
	defer rows.Close()

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	w := csv.NewWriter(file)

	if err := w.Write([]string{"file_path", "tag_name"}); err != nil {
		file.Close()
		os.Remove(outputPath)
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	for rows.Next() {
		var filePath, tagName string
		if err := rows.Scan(&filePath, &tagName); err != nil {
			w.Flush()
			file.Close()
			os.Remove(outputPath)
			return fmt.Errorf("failed to scan row: %w", err)
		}
		if err := w.Write([]string{filePath, tagName}); err != nil {
			w.Flush()
			file.Close()
			os.Remove(outputPath)
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		w.Flush()
		file.Close()
		os.Remove(outputPath)
		return fmt.Errorf("error iterating rows: %w", err)
	}

	w.Flush()
	if err := w.Error(); err != nil {
		file.Close()
		os.Remove(outputPath)
		return fmt.Errorf("failed to flush CSV writer: %w", err)
	}

	if err := file.Close(); err != nil {
		os.Remove(outputPath)
		return fmt.Errorf("failed to close output file: %w", err)
	}

	return nil
}
