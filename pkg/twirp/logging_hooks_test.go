package twirp

import (
	"context"
	"net/http"
	"testing"

	"github.com/twitchtv/twirp"

	"github.com/stretchr/testify/assert"
)

func TestLoggingHooksSet(t *testing.T) {
	sut := LoggingHooks()
	assert.NotNil(t, sut.RequestReceived)
	assert.NotNil(t, sut.RequestRouted)
	assert.NotNil(t, sut.ResponseSent)
}

func TestLoggingHooksPipe(t *testing.T) {
	ctx := context.TODO()
	sut := LoggingHooks()
	twirp.WithHTTPRequestHeaders(ctx, http.Header{})

	ctx, err := sut.RequestReceived(ctx)
	if err != nil {
		t.Error("could not call twirp request received hook")
	}

	ctx, err = sut.RequestRouted(ctx)
	if err != nil {
		t.Error("could not call twirp request routed hook")
	}

	sut.ResponseSent(ctx)
}
