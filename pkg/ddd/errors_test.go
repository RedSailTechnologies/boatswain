package ddd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDDDErrorsImplementErrorsInterface(t *testing.T) {
	assert.Error(t, DestroyedError{})
	assert.Error(t, IDError{})
	assert.Error(t, InvalidArgumentError{})
	assert.Error(t, NotFoundError{})
	assert.Error(t, RequiredArgumentError{})
	assert.Error(t, UnsupportedEventError{})

	assert.NotEmpty(t, DestroyedError{Entity: "e"}.Error())
	assert.NotEmpty(t, IDError{}.Error())
	assert.NotEmpty(t, InvalidArgumentError{Arg: "arg", Val: "validation"}.Error())
	assert.NotEmpty(t, NotFoundError{Entity: "e"}.Error())
	assert.NotEmpty(t, RequiredArgumentError{Arg: "arg"}.Error())
	assert.NotEmpty(t, UnsupportedEventError{EventType: "event", Type: "e"}.Error())
}
