package twirp

import (
	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/auth"
)

// JWTHook is the twirp server hook to validate JWTs
func JWTHook(a auth.Agent) *twirp.ServerHooks {
	hooks := &twirp.ServerHooks{}
	hooks.RequestRouted = a.Authenticate
	return hooks
}
