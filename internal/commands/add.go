package commands

import (
	"database/sql"
	"fmt"
	"strings"
)

// AddTag adds a tag to a file, creating both if they don't exist
func AddTag(db *sql.DB, filePath, tagName string) error {
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

	_, err = tx.Exec("INSERT INTO files(path) VALUES(?) ON CONFLICT(path) DO NOTHING", filePath)
	if err != nil {
		return fmt.Errorf("failed to insert file: %w", err)
	}

	var fileID int64
	err = tx.QueryRow("SELECT id FROM files WHERE path = ?", filePath).Scan(&fileID)
	if err != nil {
		return fmt.Errorf("failed to get file id: %w", err)
	}

	_, err = tx.Exec("INSERT INTO tags(name) VALUES(?) ON CONFLICT(name) DO NOTHING", tagName)
	if err != nil {
		return fmt.Errorf("failed to insert tag: %w", err)
	}

	var tagID int64
	err = tx.QueryRow("SELECT id FROM tags WHERE name = ?", tagName).Scan(&tagID)
	if err != nil {
		return fmt.Errorf("failed to get tag id: %w", err)
	}

	_, err = tx.Exec("INSERT OR IGNORE INTO file_tags(file_id, tag_id) VALUES(?, ?)", fileID, tagID)
	if err != nil {
		return fmt.Errorf("failed to link file and tag: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
