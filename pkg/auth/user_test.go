package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasRole(t *testing.T) {
	sut := &User{}
	empty := sut.hasRole("some role")
	sut.Roles = []string{
		"some role",
		"some other role",
	}

	f1 := sut.hasRole("some role")
	f2 := sut.hasRole("some other role")
	nf := sut.hasRole("some nonexistent role")

	assert.False(t, empty)
	assert.True(t, f1)
	assert.True(t, f2)
	assert.False(t, nf)
}

func TestValidScope(t *testing.T) {
	sut := &User{Scope: "scope"}

	noErr := sut.validateScope("scope")
	err := sut.validateScope("three sixty no")

	assert.Nil(t, noErr)
	assert.Error(t, err)
}

func TestUserErrors(t *testing.T) {
	assert.Error(t, userScopeError{})

	assert.NotEmpty(t, userScopeError{}.Error())
}
