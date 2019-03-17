package validate

import (
	"fmt"
)

// Error holds validate error's information
type Error struct {
	Field   string
	Message string
}

func (err *Error) Error() string {
	return fmt.Sprintf("validate: %s %s", err.Field, err.Message)
}

// NewError creates new validate error
func NewError(field, message string) error {
	return &Error{
		Field:   field,
		Message: message,
	}
}

// NewRequiredError creates new required validate error
func NewRequiredError(field string) error {
	return NewError(field, "required")
}
