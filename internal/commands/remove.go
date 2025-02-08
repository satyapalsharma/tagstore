package commands

import (
	"database/sql"
	"fmt"
	"strings"
)

// RemoveTag removes a tag from a file in the database
func RemoveTag(db *sql.DB, filePath, tagName string) error {
	filePath = strings.TrimSpace(filePath)
	tagName = strings.TrimSpace(tagName)

	if filePath == "" || tagName == "" {
		return fmt.Errorf("file path and tag name cannot be empty")
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var fileID int64
	err = tx.QueryRow("SELECT id FROM files WHERE path = ?", filePath).Scan(&fileID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("file not found: %s", filePath)
		}
		return fmt.Errorf("failed to query file: %w", err)
	}

	var tagID int64
	err = tx.QueryRow("SELECT id FROM tags WHERE name = ?", tagName).Scan(&tagID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("tag not found: %s", tagName)
		}
		return fmt.Errorf("failed to query tag: %w", err)
	}

	result, err := tx.Exec("DELETE FROM file_tags WHERE file_id = ? AND tag_id = ?", fileID, tagID)
	if err != nil {
		return fmt.Errorf("failed to remove tag: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tag '%s' is not associated with file '%s'", tagName, filePath)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
