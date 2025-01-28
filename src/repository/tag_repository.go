package repository

import (
	"context"
	"database/sql"
	"fmt"

	"tagstore/src/lib"
	"tagstore/src/models"
)

// TagRepository implements types.TagRepository using SQLite
type TagRepository struct {
	db *lib.DB
}

// NewTagRepository creates a new TagRepository
func NewTagRepository(db *lib.DB) *TagRepository {
	return &TagRepository{db: db}
}

// Create inserts a new tag into the database
func (r *TagRepository) Create(ctx context.Context, tag *models.Tag) error {
	query := `INSERT INTO tags (name) VALUES (?) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(query, tag.Name).Scan(&tag.ID, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		if isUniqueConstraintError(err) {
			return models.ErrDuplicateTag
		}
		return fmt.Errorf("failed to create tag: %w", err)
	}
	return nil
}

// GetByName retrieves a tag by its name
func (r *TagRepository) GetByName(ctx context.Context, name string) (*models.Tag, error) {
	tag := &models.Tag{}
	query := `SELECT id, name, created_at, updated_at FROM tags WHERE name = ?`
	err := r.db.QueryRow(query, name).Scan(&tag.ID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrTagNotFound
		}
		return nil, fmt.Errorf("failed to get tag by name: %w", err)
	}
	return tag, nil
}

// GetByID retrieves a tag by its ID
func (r *TagRepository) GetByID(ctx context.Context, id int64) (*models.Tag, error) {
	tag := &models.Tag{}
	query := `SELECT id, name, created_at, updated_at FROM tags WHERE id = ?`
	err := r.db.QueryRow(query, id).Scan(&tag.ID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrTagNotFound
		}
		return nil, fmt.Errorf("failed to get tag by id: %w", err)
	}
	return tag, nil
}

// List retrieves all tags
func (r *TagRepository) List(ctx context.Context) ([]models.Tag, error) {
	query := `SELECT id, name, created_at, updated_at FROM tags ORDER BY name`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
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

// Delete removes a tag by its ID
func (r *TagRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM tags WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return models.ErrTagNotFound
	}
	return nil
}

func isUniqueConstraintError(err error) bool {
	return err != nil && (err.Error() == "UNIQUE constraint failed: tags.name" ||
		err.Error() == "constraint failed")
}
