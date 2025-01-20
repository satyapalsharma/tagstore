package types

import (
	"context"
	"tagstore/src/models"
)

// TagRepository defines the interface for tag data operations
	Create(ctx context.Context, tag *models.Tag) error
	GetByName(ctx context.Context, name string) (*models.Tag, error)
	GetByID(ctx context.Context, id int64) (*models.Tag, error)
	List(ctx context.Context) ([]models.Tag, error)
	Delete(ctx context.Context, id int64) error
}

// FileRepository defines the interface for file data operations
	Create(ctx context.Context, file *models.File) error
	GetByPath(ctx context.Context, path string) (*models.File, error)
	GetByID(ctx context.Context, id int64) (*models.File, error)
	List(ctx context.Context, limit int) ([]models.File, error)
	Delete(ctx context.Context, id int64) error
}

// FileTagRepository defines the interface for file-tag relationship operations
	AddTag(ctx context.Context, fileID int64, tagID int64) error
	RemoveTag(ctx context.Context, fileID int64, tagID int64) error
	GetTagsForFile(ctx context.Context, fileID int64) ([]models.Tag, error)
	GetFilesForTag(ctx context.Context, tagID int64) ([]models.File, error)
	SearchFiles(ctx context.Context, query *TagQuery) ([]models.File, error)
	HasTag(ctx context.Context, fileID int64, tagID int64) (bool, error)
}
