package twirp

import (
	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/auth"
)

// JWTHook is the twirp server hook to validate JWTs
func JWTHook() *twirp.ServerHooks {
	hooks := &twirp.ServerHooks{}
	hooks.RequestRouted = auth.ValidateJWT
	return hooks
}
