package errorhandler

import (
	"errors"
)

type ErrorCode string

var (
	ErrInternalServer               = errors.New("internal server error")
	INVALID_INPUT         ErrorCode = "INVALID_INPUT"
	NOT_FOUND             ErrorCode = "NOT_FOUND"
	CONFLICT              ErrorCode = "CONFLICT"
	INTERNAL_SERVER_ERROR ErrorCode = "INTERNAL_SERVER_ERROR"
)

// Custom Errors
type CustomError struct {
	Message string    `json:"message"`
	Code    ErrorCode `json:"code"` // HTTP status code
}

func (e *CustomError) Error() string {
	return e.Message
}

// NewCustomError creates a new CustomError instance
func NewCustomError(message string, code ErrorCode) *CustomError {
	return &CustomError{Message: message, Code: code}
}

// NotFoundError is used when a resource is not found.
type NotFoundError struct {
	Message string `json:"message"`
}

func (e *NotFoundError) Error() string {
	return e.Message
}

// NewNotFoundError creates a new NotFoundError.
func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{Message: message}
}

// InvalidInputError is used for invalid input data.
type InvalidInputError struct {
	Message string `json:"message"`
}

func (e *InvalidInputError) Error() string {
	return e.Message
}

// NewInvalidInputError creates a new InvalidInputError.
func NewInvalidInputError(message string) *InvalidInputError {
	return &InvalidInputError{Message: message}
}

// InternalServerError represents a generic internal server error.
type InternalServerError struct {
	Message string `json:"message"`
}

func (e *InternalServerError) Error() string {
	return e.Message
}

// NewInternalServerError creates a new InternalServerError.
func NewInternalServerError(message string) *InternalServerError {
	return &InternalServerError{Message: message}
}

type SightingTooCloseError struct {
	Message string `json:"message"`
}

func (e *SightingTooCloseError) Error() string {
	return e.Message
}

func As(err error, target any) bool {
	return errors.As(err, target)
}
