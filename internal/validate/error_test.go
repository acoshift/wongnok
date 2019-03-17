package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	// 1. use syntax error
	var _ error = &Error{}

	// 2. run test
	assert.Implements(t, (*error)(nil), &Error{})
}
