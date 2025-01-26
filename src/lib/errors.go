package lib

import (
	"errors"
	"fmt"
)

// ErrorType represents the category of error
type ErrorType string

const (
	ErrorTypeBadRequest   ErrorType = "bad_request"
	ErrorTypeUnauthorized ErrorType = "unauthorized"
	ErrorTypeForbidden    ErrorType = "forbidden"
	ErrorTypeNotFound     ErrorType = "not_found"
	ErrorTypeConflict     ErrorType = "conflict"
	ErrorTypeInternal     ErrorType = "internal"
)

// AppError represents a structured application error
type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
	}

func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new AppError
func NewAppError(errType ErrorType, message string, err error) *AppError {
	return &AppError{
		Type:    errType,
		Message: message,
		Err:     err,
	}
}

// BadRequestError creates a bad request error
func BadRequestError(message string, err error) *AppError {
	return NewAppError(ErrorTypeBadRequest, message, err)
}

// UnauthorizedError creates an unauthorized error
func UnauthorizedError(message string, err error) *AppError {
	return NewAppError(ErrorTypeUnauthorized, message, err)
}

// NotFoundError creates a not found error
func NotFoundError(message string, err error) *AppError {
	return NewAppError(ErrorTypeNotFound, message, err)
}

// ConflictError creates a conflict error
func ConflictError(message string, err error) *AppError {
	return NewAppError(ErrorTypeConflict, message, err)
}

// InternalError creates an internal server error
func InternalError(message string, err error) *AppError {
	return NewAppError(ErrorTypeInternal, message, err)
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// GetAppErrorType extracts the error type from an error
func GetAppErrorType(err error) ErrorType {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type
	}
	return ErrorTypeInternal
}
