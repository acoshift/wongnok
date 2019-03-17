package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateError(t *testing.T) {
	// 1. use syntax error
	var _ error = &ValidateError{}

	// 2. run test
	assert.Implements(t, (*error)(nil), &ValidateError{})
}
