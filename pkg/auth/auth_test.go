package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthErrors(t *testing.T) {
	assert.Error(t, AuthenticationError{})
	assert.Error(t, NotAuthorizedError{})
	assert.Error(t, userContextError{})

	assert.NotEmpty(t, AuthenticationError{}.Error())
	assert.NotEmpty(t, NotAuthorizedError{}.Error())
	assert.NotEmpty(t, userContextError{}.Error())
}
