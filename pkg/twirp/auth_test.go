package twirp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateJWTAssignedToRouted(t *testing.T) {
	sut := JWTHook()
	assert.True(t, sut.RequestRouted != nil)
}
