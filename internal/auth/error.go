package auth

import (
	"errors"
)

// Errors
var (
	ErrUsernameRequired     = errors.New("auth: username required")
	ErrUsernameTooShort     = errors.New("auth: username too short")
	ErrUsernameTooLong      = errors.New("auth: username too long")
	ErrUsernameInvalid      = errors.New("auth: username invalid")
	ErrUsernameNotAvailable = errors.New("auth: username not available")
)
