package repository

import (
	"context"
	"fmt"
	"strings"

	"tagstore/src/lib"
	"tagstore/src/models"
	"tagstore/src/types"
)

// FileTagRepository implements types.FileTagRepository using SQLite
type FileTagRepository struct {
	db *lib.DB
}

// NewFileTagRepository creates a new FileTagRepository
func NewFileTagRepository(db *lib.DB) *FileTagRepository {
	return &FileTagRepository{db: db}
}

// AddTag associates a tag with a file
func (r *FileTagRepository) AddTag(ctx context.Context, fileID int64, tagID int64) error {
	query := `INSERT INTO file_tags (file_id, tag_id) VALUES (?, ?)`
	_, err := r.db.Exec(query, fileID, tagID)
	if err != nil {
		if isUniqueConstraintError(err) {
			return nil // Already exists, not an error
		}
		return fmt.Errorf("failed to add tag to file: %w", err)
	}
	return nil
}

// RemoveTag disassociates a tag from a file
func (r *FileTagRepository) RemoveTag(ctx context.Context, fileID int64, tagID int64) error {
	query := `DELETE FROM file_tags WHERE file_id = ? AND tag_id = ?`
	result, err := r.db.Exec(query, fileID, tagID)
	if err != nil {
		return fmt.Errorf("failed to remove tag from file: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return models.ErrFileNotTagged
	}
	return nil
}

// GetTagsForFile retrieves all tags associated with a file
func (r *FileTagRepository) GetTagsForFile(ctx context.Context, fileID int64) ([]models.Tag, error) {
	query := `
		SELECT t.id, t.name, t.created_at, t.updated_at
		FROM tags t
		INNER JOIN file_tags ft ON t.id = ft.tag_id
		WHERE ft.file_id = ?
		ORDER BY t.name
	`
	rows, err := r.db.Query(query, fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags for file: %w", err)
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var tag models.Tag
		if err := rows.Scan(&tag.ID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// GetFilesForTag retrieves all files associated with a tag
func (r *FileTagRepository) GetFilesForTag(ctx context.Context, tagID int64) ([]models.File, error) {
	query := `
		SELECT f.id, f.path, f.hash, f.created_at, f.updated_at
		FROM files f
		INNER JOIN file_tags ft ON f.id = ft.file_id
		WHERE ft.tag_id = ?
		ORDER BY f.path
	`
	rows, err := r.db.Query(query, tagID)
	if err != nil {
		return nil, fmt.Errorf("failed to get files for tag: %w", err)
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var file models.File
		if err := rows.Scan(&file.ID, &file.Path, &file.Hash, &file.CreatedAt, &file.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}
		files = append(files, file)
	}

	return files, nil
}

// HasTag checks if a file has a specific tag
func (r *FileTagRepository) HasTag(ctx context.Context, fileID int64, tagID int64) (bool, error) {
	query := `SELECT COUNT(*) FROM file_tags WHERE file_id = ? AND tag_id = ?`
	var count int
	err := r.db.QueryRow(query, fileID, tagID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check tag association: %w", err)
	}
	return count > 0, nil
}

// SearchFiles finds files matching a boolean tag query
func (r *FileTagRepository) SearchFiles(ctx context.Context, query *types.TagQuery) ([]models.File, error) {
	if query == nil || query.Raw == "" {
		return r.getAllFiles(ctx)
	}

	// Build dynamic SQL based on query
	var conditions []string
	var args []interface{}

	// Positive tags (AND) - file must have ALL of these
	for _, tag := range query.Positive {
		conditions = append(conditions, `
			EXISTS (
				SELECT 1 FROM file_tags ft2
				INNER JOIN tags t2 ON ft2.tag_id = t2.id
				WHERE ft2.file_id = f.id AND t2.name = ?
			)
		`)
		args = append(args, tag)
	}

	// Negative tags (NOT) - file must NOT have ANY of these
	for _, tag := range query.Negative {
		conditions = append(conditions, `
			NOT EXISTS (
				SELECT 1 FROM file_tags ft3
				INNER JOIN tags t3 ON ft3.tag_id = t3.id
				WHERE ft3.file_id = f.id AND t3.name = ?
			)
		`)
		args = append(args, tag)
	}

	// Optional tags (OR) - file must have AT LEAST ONE of these
	if len(query.Optional) > 0 {
		optionalConditions := make([]string, len(query.Optional))
		for i, tag := range query.Optional {
			optionalConditions[i] = `
				EXISTS (
					SELECT 1 FROM file_tags ft4
					INNER JOIN tags t4 ON ft4.tag_id = t4.id
					WHERE ft4.file_id = f.id AND t4.name = ?
				)
			`
			args = append(args, tag)
		}
		conditions = append(conditions, "("+strings.Join(optionalConditions, " OR ")+")")
	}

	whereClause := strings.Join(conditions, " AND ")
	if whereClause == "" {
		whereClause = "1=1"
	}

	queryStr := fmt.Sprintf(`
		SELECT DISTINCT f.id, f.path, f.hash, f.created_at, f.updated_at
		FROM files f
		WHERE %s
		ORDER BY f.path
	`, whereClause)

	rows, err := r.db.Query(queryStr, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search files: %w", err)
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var file models.File
		if err := rows.Scan(&file.ID, &file.Path, &file.Hash, &file.CreatedAt, &file.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}
		files = append(files, file)
	}

	return files, nil
}

// SearchFilesByTags finds files that have all or any of the specified tags
func (r *FileTagRepository) SearchFilesByTags(ctx context.Context, tags []string, matchAll bool) ([]models.File, error) {
	if len(tags) == 0 {
		return nil, fmt.Errorf("at least one tag is required")
	}

	// Trim tags and validate
	trimmedTags := make([]string, len(tags))
	for i, tag := range tags {
		trimmed := strings.TrimSpace(tag)
		if trimmed == "" {
			return nil, fmt.Errorf("tag cannot be empty")
		}
		trimmedTags[i] = trimmed
	}

	var query string
	var args []interface{}

	placeholders := make([]string, len(trimmedTags))
	for i := range trimmedTags {
		placeholders[i] = "?"
	}
	inClause := strings.Join(placeholders, ",")

	if matchAll {
		// AND logic: file must have ALL tags
		query = fmt.Sprintf(`
			SELECT f.id, f.path, f.hash, f.created_at, f.updated_at
			FROM files f
			JOIN file_tags ft ON f.id = ft.file_id
			JOIN tags t ON t.id = ft.tag_id
			WHERE t.name IN (%s)
			GROUP BY f.id
			HAVING COUNT(DISTINCT t.id) = ?
		`, inClause)

		args = make([]interface{}, len(trimmedTags)+1)
		for i, tag := range trimmedTags {
			args[i] = tag
		}
		args[len(trimmedTags)] = len(trimmedTags)
	} else {
		// OR logic: file must have ANY tag
		query = fmt.Sprintf(`
			SELECT DISTINCT f.id, f.path, f.hash, f.created_at, f.updated_at
			FROM files f
			JOIN file_tags ft ON f.id = ft.file_id
			JOIN tags t ON t.id = ft.tag_id
			WHERE t.name IN (%s)
			ORDER BY f.path
		`, inClause)

		args = make([]interface{}, len(trimmedTags))
		for i, tag := range trimmedTags {
			args[i] = tag
		}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search files by tags: %w", err)
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var file models.File
		if err := rows.Scan(&file.ID, &file.Path, &file.Hash, &file.CreatedAt, &file.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}
		files = append(files, file)
	}

	return files, nil
}

func (r *FileTagRepository) getAllFiles(ctx context.Context) ([]models.File, error) {
	query := `SELECT id, path, hash, created_at, updated_at FROM files ORDER BY path`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all files: %w", err)
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var file models.File
		if err := rows.Scan(&file.ID, &file.Path, &file.Hash, &file.CreatedAt, &file.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}
		files = append(files, file)
	}

	return files, nil
}

straintError(err error) bool {
	return err != nil && (err.Error() == "UNIQUE constraint failed: file_tags.file_id, file_tags.tag_id" ||
		err.Error() == "constraint failed")
}
