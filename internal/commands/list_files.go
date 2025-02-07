package commands

import (
	"database/sql"
	"fmt"
	"strings"
)

// ListFilesForTag returns all file paths associated with the given tag name.
func ListFilesForTag(db *sql.DB, tagName string) ([]string, error) {
	tagName = strings.TrimSpace(tagName)
	if tagName == "" {
		return nil, fmt.Errorf("tag name cannot be empty")
	}

	var tagID int64
	err := db.QueryRow("SELECT id FROM tags WHERE name = ?", tagName).Scan(&tagID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tag not found: %s", tagName)
		}
		return nil, fmt.Errorf("failed to query tag: %w", err)
	}

	rows, err := db.Query("SELECT f.path FROM files f JOIN file_tags ft ON f.id = ft.file_id WHERE ft.tag_id = ? ORDER BY f.path ASC", tagID)
	if err != nil {
		return nil, fmt.Errorf("failed to query files: %w", err)
	}
	defer rows.Close()

	var paths []string
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return nil, fmt.Errorf("failed to scan file path: %w", err)
		}
		paths = append(paths, path)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return paths, nil
}
