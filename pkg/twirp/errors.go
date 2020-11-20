package twirp

import (
	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/twitchtv/twirp"
)

// ToTwirpError converts an error to a Twirp response error
func ToTwirpError(e error, m string) error {
	switch e.(type) {
	case ddd.DestroyedError:
		return twirp.NotFoundError(e.Error())
	case ddd.InvalidArgumentError:
		return twirp.InvalidArgumentError(e.(ddd.InvalidArgumentError).Arg, e.Error())
	case ddd.NotFoundError:
		return twirp.NotFoundError(e.Error())
	case ddd.RequiredArgumentError:
		return twirp.RequiredArgumentError(e.(ddd.RequiredArgumentError).Arg)
	default:
		return twirp.InternalError(m)
	}
}
