package models

import "errors"

// Common domain errors
var (
	ErrEmptyTagName    = errors.New("tag name cannot be empty")
	ErrTagNameTooLong  = errors.New("tag name exceeds maximum length of 100 characters")
	ErrEmptyFilePath   = errors.New("file path cannot be empty")
	ErrFileNotFound    = errors.New("file not found")
	ErrTagNotFound     = errors.New("tag not found")
	ErrDuplicateTag    = errors.New("tag already exists")
	ErrDuplicateFile   = errors.New("file already exists")
	ErrFileNotTagged   = errors.New("file is not tagged with the specified tag")
	ErrQueryParseError = errors.New("invalid query syntax")
)
