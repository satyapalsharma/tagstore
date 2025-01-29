package utils

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// ValidateFilePath checks if a file path is valid
func ValidateFilePath(path string) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	// Check for absolute path
	if !filepath.IsAbs(path) {
		return fmt.Errorf("file path must be absolute: %s", path)
	}

	// Check for path traversal attempts
	cleaned := filepath.Clean(path)
	if strings.Contains(cleaned, "..") {
		return fmt.Errorf("invalid file path: %s", path)
	}

	return nil
}

// ValidateTagName checks if a tag name is valid
func ValidateTagName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("tag name cannot be empty")
	}

	if len(name) > 100 {
		return fmt.Errorf("tag name exceeds maximum length of 100 characters")
	}

	// Tag names should be alphanumeric with hyphens and underscores
	validTag := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validTag.MatchString(name) {
		return fmt.Errorf("tag name contains invalid characters: %s", name)
	}

	return nil
}
