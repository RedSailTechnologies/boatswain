package twirp

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/logger"
)

var timeKey = new(int)
var idKey = new(int)

// LoggingHooks creates server hooks which log on requests and responses automatically
func LoggingHooks() *twirp.ServerHooks {
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
			logger.Info("request received",
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
			logger.Info("response sent",
				"id", id,
				"service", service,
				"method", method,
				"status", status,
				"elapsed", time.Since(start))
		},
	}
}

// LoggingException represents a service method we don't want to log
type LoggingException struct {
	Method  string
	Service string
}

// LoggingHooksWithExceptions creates server hooks which log on requests and responses automatically
func LoggingHooksWithExceptions(exceptions ...LoggingException) *twirp.ServerHooks {
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
			if !loggingExceptions(exceptions).isException(method, service) {
				logger.Info("request received",
					"id", id,
					"service", service,
					"method", method,
					"headers", headers)
			}
			return ctx, nil
		},
		ResponseSent: func(ctx context.Context) {
			start := ctx.Value(timeKey).(time.Time)
			id := ctx.Value(idKey).(uuid.UUID)
			method, _ := twirp.MethodName(ctx)
			service, _ := twirp.ServiceName(ctx)
			status, _ := twirp.StatusCode(ctx)
			if !loggingExceptions(exceptions).isException(method, service) {
				logger.Info("response sent",
					"id", id,
					"service", service,
					"method", method,
					"status", status,
					"elapsed", time.Since(start))
			}
		},
	}
}

type loggingExceptions []LoggingException

func (e loggingExceptions) isException(method, service string) bool {
	for _, ex := range e {
		if ex.Service == service && ex.Method == method {
			return true
		}
	}
	return false
}
