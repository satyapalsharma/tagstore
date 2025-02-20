package commands

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

// ImportFromCSV reads a CSV file and imports file-tag associations into the database.
// The CSV must have a header row with exactly two columns: file_path and tag_name.
// Returns the number of new associations created.
func ImportFromCSV(db *sql.DB, inputPath string) (int, error) {
	if inputPath == "" {
		return 0, fmt.Errorf("input path cannot be empty")
	}

	f, err := os.Open(inputPath)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return 0, fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) == 0 {
		return 0, fmt.Errorf("CSV file is empty")
	}

	header := records[0]
	if len(header) != 2 || header[0] != "file_path" || header[1] != "tag_name" {
		return 0, fmt.Errorf("invalid CSV header: expected [file_path, tag_name]")
	}

	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	count := 0
	for i, record := range records[1:] {
		line := i + 2 // 1-indexed line number, plus header row
		if len(record) != 2 {
			tx.Rollback()
			return 0, fmt.Errorf("invalid record at line %d: expected 2 fields, got %d", line, len(record))
		}

		filePath := strings.TrimSpace(record[0])
		tagName := strings.TrimSpace(record[1])

		if filePath == "" || tagName == "" {
			tx.Rollback()
			return 0, fmt.Errorf("empty field at line %d", line)
		}

		// Insert file if not exists, get file ID
		var fileID int64
		err = tx.QueryRow("SELECT id FROM files WHERE path = ?", filePath).Scan(&fileID)
		if err != nil {
			if err == sql.ErrNoRows {
				err = tx.QueryRow("INSERT INTO files (path) VALUES (?) RETURNING id", filePath).Scan(&fileID)
				if err != nil {
					tx.Rollback()
					return 0, fmt.Errorf("failed to insert file at line %d: %w", line, err)
				}
			} else {
				tx.Rollback()
				return 0, fmt.Errorf("failed to get file id at line %d: %w", line, err)
			}
		}

		// Insert tag if not exists, get tag ID
		var tagID int64
		err = tx.QueryRow("SELECT id FROM tags WHERE name = ?", tagName).Scan(&tagID)
		if err != nil {
			if err == sql.ErrNoRows {
				err = tx.QueryRow("INSERT INTO tags (name) VALUES (?) RETURNING id", tagName).Scan(&tagID)
				if err != nil {
					tx.Rollback()
					return 0, fmt.Errorf("failed to insert tag at line %d: %w", line, err)
				}
			} else {
				tx.Rollback()
				return 0, fmt.Errorf("failed to get tag id at line %d: %w", line, err)
			}
		}

		// Insert association
		result, err := tx.Exec("INSERT OR IGNORE INTO file_tags (file_id, tag_id) VALUES (?, ?)", fileID, tagID)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert association at line %d: %w", line, err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to get rows affected at line %d: %w", line, err)
		}

		if rowsAffected > 0 {
			count++
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return count, nil
}
