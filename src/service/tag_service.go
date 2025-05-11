package service

import (
	"context"

	"tagstore/src/lib"
	"tagstore/src/models"
	"tagstore/src/repository"
)

// TagService provides business logic for tag operations
type TagService struct {
	tagRepo   *repository.TagRepository
	fileRepo  *repository.FileRepository
	fileTagRepo *repository.FileTagRepository
	logger    *lib.Logger
}

// NewTagService creates a new TagService
func NewTagService(
	tagRepo *repository.TagRepository,
	fileRepo *repository.FileRepository,
	fileTagRepo *repository.FileTagRepository,
	logger *lib.Logger,
) *TagService {
	return &TagService{
		tagRepo:    tagRepo,
		fileRepo:   fileRepo,
		fileTagRepo: fileTagRepo,
		logger:     logger,
	}
}

// AddTagToFile adds a tag to a file, creating both if they don't exist
func (s *TagService) AddTagToFile(ctx context.Context, filePath string, tagName string) error {
	// Validate inputs
	if filePath == "" {
		return lib.BadRequestError("file path cannot be empty", nil)
	}
	if tagName == "" {
		return lib.BadRequestError("tag name cannot be empty", nil)
	}

	// Get or create file
	file, err := s.fileRepo.GetByPath(ctx, filePath)
	if err != nil {
		if err == models.ErrFileNotFound {
			file = &models.File{Path: filePath}
			if err := file.Validate(); err != nil {
				return lib.BadRequestError(err.Error(), err)
			}
			if err := s.fileRepo.Create(ctx, file); err != nil {
				s.logger.Error("Failed to create file: %v", err)
				return lib.InternalError("failed to create file", err)
			}
		} else {
			s.logger.Error("Failed to get file: %v", err)
			return lib.InternalError("failed to get file", err)
		}
	}

	// Get or create tag
	tag, err := s.tagRepo.GetByName(ctx, tagName)
	if err != nil {
		if err == models.ErrTagNotFound {
			tag = &models.Tag{Name: tagName}
			if err := tag.Validate(); err != nil {
				return lib.BadRequestError(err.Error(), err)
			}
			if err := s.tagRepo.Create(ctx, tag); err != nil {
				s.logger.Error("Failed to create tag: %v", err)
				return lib.InternalError("failed to create tag", err)
			}
		} else {
			s.logger.Error("Failed to get tag: %v", err)
			return lib.InternalError("failed to get tag", err)
		}
	}

	// Associate tag with file
	if err := s.fileTagRepo.AddTag(ctx, file.ID, tag.ID); err != nil {
		s.logger.Error("Failed to add tag to file: %v", err)
		return lib.InternalError("failed to add tag to file", err)
	}

	s.logger.Info("Added tag '%s' to file '%s'", tagName, filePath)
	return nil
}

// RemoveTagFromFile removes a tag from a file
func (s *TagService) RemoveTagFromFile(ctx context.Context, filePath string, tagName string) error {
	if filePath == "" {
		return lib.BadRequestError("file path cannot be empty", nil)
	}
	if tagName == "" {
		return lib.BadRequestError("tag name cannot be empty", nil)
	}

	file, err := s.fileRepo.GetByPath(ctx, filePath)
	if err != nil {
		if err == models.ErrFileNotFound {
			return lib.NotFoundError("file not found", err)
		}
		s.logger.Error("Failed to get file: %v", err)
		return lib.InternalError("failed to get file", err)
	}

	tag, err := s.tagRepo.GetByName(ctx, tagName)
	if err != nil {
		if err == models.ErrTagNotFound {
			return lib.NotFoundError("tag not found", err)
		}
		s.logger.Error("Failed to get tag: %v", err)
		return lib.InternalError("failed to get tag", err)
	}

	if err := s.fileTagRepo.RemoveTag(ctx, file.ID, tag.ID); err != nil {
		if err == models.ErrFileNotTagged {
			return lib.NotFoundError("file is not tagged with the specified tag", err)
		}
		s.logger.Error("Failed to remove tag from file: %v", err)
		return lib.InternalError("failed to remove tag from file", err)
	}

	s.logger.Info("Removed tag '%s' from file '%s'", tagName, filePath)
	return nil
}

// ListTagsForFile retrieves all tags for a file
func (s *TagService) ListTagsForFile(ctx context.Context, filePath string) ([]models.Tag, error) {
	if filePath == "" {
		return nil, lib.BadRequestError("file path cannot be empty", nil)
	}

	file, err := s.fileRepo.GetByPath(ctx, filePath)
	if err != nil {
		if err == models.ErrFileNotFound {
			return nil, lib.NotFoundError("file not found", err)
		}
		s.logger.Error("Failed to get file: %v", err)
		return nil, lib.InternalError("failed to get file", err)
	}

	tags, err := s.fileTagRepo.GetTagsForFile(ctx, file.ID)
	if err != nil {
		s.logger.Error("Failed to get tags for file: %v", err)
		return nil, lib.InternalError("failed to get tags for file", err)
	}

	return tags, nil
}

// ListFilesForTag retrieves all files for a tag
func (s *TagService) ListFilesForTag(ctx context.Context, tagName string) ([]models.File, error) {
	if tagName == "" {
		return nil, lib.BadRequestError("tag name cannot be empty", nil)
	}

	tag, err := s.tagRepo.GetByName(ctx, tagName)
	if err != nil {
		if err == models.ErrTagNotFound {
			return nil, lib.NotFoundError("tag not found", err)
		}
		s.logger.Error("Failed to get tag: %v", err)
		return nil, lib.InternalError("failed to get tag", err)
	}

	files, err := s.fileTagRepo.GetFilesForTag(ctx, tag.ID)
	if err != nil {
		s.logger.Error("Failed to get files for tag: %v", err)
		return nil, lib.InternalError("failed to get files for tag", err)
	}

	return files, nil
}
