package logger

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/twitchtv/twirp"
)

var timeKey = new(int)
var idKey = new(int)

// TwirpHooks creates server hooks which log on requests and responses automatically
func TwirpHooks() *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestReceived: func(ctx context.Context) (context.Context, error) {
			start := time.Now()
			id := uuid.New()
			ctx = context.WithValue(ctx, timeKey, start)
			ctx = context.WithValue(ctx, idKey, id)
			return ctx, nil
		},
		RequestRouted: func(ctx context.Context) (context.Context, error) {
			id := ctx.Value(idKey).(uuid.UUID)
			headers, _ := twirp.HTTPRequestHeaders(ctx)
			method, _ := twirp.MethodName(ctx)
			service, _ := twirp.ServiceName(ctx)
			l.Infow("request received",
				"id", id,
				"service", service,
				"method", method,
				"headers", headers)
			return ctx, nil
		},
		ResponseSent: func(ctx context.Context) {
			start := ctx.Value(timeKey).(time.Time)
			id := ctx.Value(idKey).(uuid.UUID)
			method, _ := twirp.MethodName(ctx)
			service, _ := twirp.ServiceName(ctx)
			status, _ := twirp.StatusCode(ctx)
			l.Infow("response sent",
				"id", id,
				"service", service,
				"method", method,
				"status", status,
				"elapsed", time.Since(start))
		},
	}
}
