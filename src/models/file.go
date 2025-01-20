package models

import (
	"time"
)

// File represents a file entity in the database
type File struct {
	ID        int64     `db:"id" json:"id"`
	Path      string    `db:"path" json:"path"`
	Hash      string    `db:"hash" json:"hash"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// TableName returns the database table name for File
func (f File) TableName() string {
	return "files"
}

// Validate checks if the File model is valid
func (f File) Validate() error {
	if strings.TrimSpace(f.Path) == "" {
		return ErrEmptyFilePath
	}
	return nil
}
