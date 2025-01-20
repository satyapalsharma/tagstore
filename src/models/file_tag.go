package models

import "time"

// FileTag represents the join table between files and tags
type FileTag struct {
	FileID    int64     `db:"file_id" json:"file_id"`
	TagID     int64     `db:"tag_id" json:"tag_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// TableName returns the database table name for FileTag
func (ft FileTag) TableName() string {
	return "file_tags"
}
