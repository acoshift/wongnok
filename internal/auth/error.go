package auth

import (
	"errors"
	"fmt"
)

// Errors
var (
	ErrUsernameRequired     = errors.New("auth: username required")
	ErrUsernameTooShort     = errors.New("auth: username too short")
	ErrUsernameTooLong      = errors.New("auth: username too long")
	ErrUsernameInvalid      = errors.New("auth: username invalid")
	ErrUsernameNotAvailable = errors.New("auth: username not available")
)

// ValidateError holds validate error's information
type ValidateError struct {
	Field   string
	Message string
}

func (err *ValidateError) Error() string {
	return fmt.Sprintf("validate: %s %s", err.Field, err.Message)
}

func newValidateError(field, message string) error {
	return &ValidateError{
		Field:   field,
		Message: message,
	}
}

func newRequiredError(field string) error {
	return &ValidateError{
		Field:   field,
		Message: "required",
	}
}
