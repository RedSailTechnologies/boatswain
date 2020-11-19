package ddd

import (
	"testing"

	"gotest.tools/assert"
)

func TestRequiredArgumentErrorMessage(t *testing.T) {
	assert.Equal(t, "all fields are required for a valid Cluster", RequiredArgumentError.Error())
}

func TestDestroyedErrorMessage(t *testing.T) {
	assert.Equal(t, "Cluster cannot be modified further as it has been destroyed", DestroyedError.Error())
}
