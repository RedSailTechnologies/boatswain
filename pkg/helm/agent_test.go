package helm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFullName(t *testing.T) {
	assert.Equal(t, "name-version.tgz", getFullName("name", "version"))
}
