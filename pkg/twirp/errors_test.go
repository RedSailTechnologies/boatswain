package twirp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/ddd"
)

func TestToTwirpError(t *testing.T) {
	destroyed := ddd.DestroyedError{Entity: "Anything"}
	invalid := ddd.InvalidArgumentError{Arg: "arg", Val: "validation"}
	notFound := ddd.NotFoundError{Entity: "Anything"}
	required := ddd.RequiredArgumentError{Arg: "arg"}

	assert.Equal(t, true, errsEqual(ToTwirpError(destroyed, "").(twirp.Error), twirp.NotFoundError(destroyed.Error())))
	assert.Equal(t, true, errsEqual(ToTwirpError(destroyed, "").(twirp.Error), twirp.NotFoundError(destroyed.Error())))
	assert.Equal(t, true, errsEqual(ToTwirpError(invalid, "").(twirp.Error), twirp.InvalidArgumentError(invalid.Arg, invalid.Error())))
	assert.Equal(t, true, errsEqual(ToTwirpError(notFound, "").(twirp.Error), twirp.NotFoundError(notFound.Error())))
	assert.Equal(t, true, errsEqual(ToTwirpError(required, "").(twirp.Error), twirp.RequiredArgumentError(required.Arg)))
	assert.Equal(t, true, errsEqual(ToTwirpError(errors.New(""), "message").(twirp.Error), twirp.InternalError("message")))
}

func errsEqual(err1, err2 twirp.Error) bool {
	return err1.Code() == err2.Code() &&
		err1.Msg() == err2.Msg()
}
