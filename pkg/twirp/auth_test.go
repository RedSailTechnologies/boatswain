package twirp

import (
	"testing"

	"github.com/redsailtechnologies/boatswain/pkg/auth"
	"github.com/stretchr/testify/assert"
)

func TestValidateJWTAssignedToRouted(t *testing.T) {
	sut := JWTHook(&auth.OIDCAgent{})
	assert.True(t, sut.RequestRouted != nil)
}
