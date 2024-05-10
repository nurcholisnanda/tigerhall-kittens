package errors

import (
	"errors"
	"fmt"
)

type ErrorCode string

var (
	ErrInternalServer           = errors.New("internal server error")
	INVALID_INPUT     ErrorCode = "INVALID_INPUT"
)

// Custom Errors
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

func As(err error, target any) bool {
	return errors.As(err, target)
}
