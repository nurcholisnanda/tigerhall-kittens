package errorhandler

import (
	"errors"
	"fmt"
)

type ErrorCode string

var (
	ErrInternalServer           = errors.New("internal server error")
	INVALID_INPUT     ErrorCode = "INVALID_INPUT"
	NOT_FOUND         ErrorCode = "NOT_FOUND"
	CONFLICT          ErrorCode = "CONFLICT"
)

// Custom Errors
type CustomError struct {
	Message string `json:"message"`
	Code    int    `json:"code"` // HTTP status code
}

func (e *CustomError) Error() string {
	return e.Message
}

// NewCustomError creates a new CustomError instance
func NewCustomError(message string, code int) *CustomError {
	return &CustomError{Message: message, Code: code}
}

type InvalidCoordinatesError struct {
	Message string `json:"message"`
}

func (e *InvalidCoordinatesError) Error() string {
	return e.Message
}

type InvalidDateOfBirthError struct {
	Message string `json:"message"`
}

func (e *InvalidDateOfBirthError) Error() string {
	return e.Message
}

type InvalidLastSeenTimeError struct {
	Message string `json:"message"`
}

func (e *InvalidLastSeenTimeError) Error() string {
	return e.Message
}

type TigerCreationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *TigerCreationError) Error() string {
	return fmt.Sprintf("tiger creation failed: %s: %s", e.Field, e.Message)
}

type TigerNotFound struct {
	Message string `json:"message"`
}

func (e *TigerNotFound) Error() string {
	return e.Message
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
