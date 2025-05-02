package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"tagstore/src/lib"
	"tagstore/src/models"
)

// FileRepository implements types.FileRepository using SQLite
type FileRepository struct {
	db *lib.DB
}

// NewFileRepository creates a new FileRepository
func NewFileRepository(db *lib.DB) *FileRepository {
	return &FileRepository{db: db}
}

// Create inserts a new file into the database
func (r *FileRepository) Create(ctx context.Context, file *models.File) error {
	query := `INSERT INTO files (path, hash) VALUES (?, ?) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(query, file.Path, file.Hash).Scan(&file.ID, &file.CreatedAt, &file.UpdatedAt)
	if err != nil {
		if isUniqueConstraintError(err) {
			return models.ErrDuplicateFile
		}
		return fmt.Errorf("failed to create file: %w", err)
	}
	return nil
}

// GetByPath retrieves a file by its path
func (r *FileRepository) GetByPath(ctx context.Context, path string) (*models.File, error) {
	file := &models.File{}
	query := `SELECT id, path, hash, created_at, updated_at FROM files WHERE path = ?`
	err := r.db.QueryRow(query, path).Scan(&file.ID, &file.Path, &file.Hash, &file.CreatedAt, &file.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to get file by path: %w", err)
	}
	return file, nil
}

// GetByID retrieves a file by its ID
func (r *FileRepository) GetByID(ctx context.Context, id int64) (*models.File, error) {
	file := &models.File{}
	query := `SELECT id, path, hash, created_at, updated_at FROM files WHERE id = ?`
	err := r.db.QueryRow(query, id).Scan(&file.ID, &file.Path, &file.Hash, &file.CreatedAt, &file.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to get file by id: %w", err)
	}
	return file, nil
}

// List retrieves all files with optional limit
func (r *FileRepository) List(ctx context.Context, limit int) ([]models.File, error) {
	query := `SELECT id, path, hash, created_at, updated_at FROM files ORDER BY path`
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
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

// Delete removes a file by its ID
func (r *FileRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM files WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return models.ErrFileNotFound
	}
	return nil
}

func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}

	message := err.Error()
	return strings.Contains(message, "UNIQUE constraint failed:") || message == "constraint failed"
}
