package service

import (
	"context"
	"fmt"

	"tagstore/src/lib"
	"tagstore/src/repository"
	"tagstore/src/types"
)

// SearchService provides business logic for search operations
type SearchService struct {
	fileTagRepo *repository.FileTagRepository
	logger      *lib.Logger
}

// NewSearchService creates a new SearchService
func NewSearchService(
	fileTagRepo *repository.FileTagRepository,
	logger *lib.Logger,
) *SearchService {
	return &SearchService{
		fileTagRepo: fileTagRepo,
		logger:      logger,
	}
}

// SearchFiles searches for files matching a boolean tag query
func (s *SearchService) SearchFiles(ctx context.Context, queryStr string) ([]string, error) {
	if queryStr == "" {
		return nil, lib.BadRequestError("query cannot be empty", nil)
	}

	query := types.NewTagQuery(queryStr)
	if query == nil {
		return nil, lib.BadRequestError("invalid query syntax", nil)
	}

	s.logger.Info("Searching files with query: %s", queryStr)

	files, err := s.fileTagRepo.SearchFiles(ctx, query)
	if err != nil {
		s.logger.Error("Search failed: %v", err)
		return nil, lib.InternalError("search failed", err)
	}

	paths := make([]string, len(files))
	for i, file := range files {
		paths[i] = file.Path
	}

	s.logger.Info("Found %d files matching query", len(paths))
	return paths, nil
}

// ListAllTags retrieves all tags in the system
func (s *SearchService) ListAllTags(ctx context.Context) ([]string, error) {
	// This will be implemented when we add a tag listing method to the repository
	// For now, return empty
	return nil, fmt.Errorf("not implemented")
}
