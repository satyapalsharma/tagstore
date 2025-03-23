package models

import (
	"strings"
	"time"
)

// Tag represents a tag entity in the database
type Tag struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// TableName returns the database table name for Tag
func (t Tag) TableName() string {
	return "tags"
}

// Validate checks if the Tag model is valid
func (t Tag) Validate() error {
	if strings.TrimSpace(t.Name) == "" {
		return ErrEmptyTagName
	}
	if len(t.Name) > 100 {
		return ErrTagNameTooLong
	}
	return nil
}
