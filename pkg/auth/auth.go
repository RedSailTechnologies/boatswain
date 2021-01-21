package auth

import (
	"context"
	"net/http"
)

// Agent is the interface through which most of the auth package is utilized
type Agent interface {
	Authenticate(context.Context) (context.Context, error)
	Authorize(context.Context, Role) error
	NewContext(ctx context.Context) (context.Context, error)
	User(context.Context) User
	Wrap(http.Handler) http.Handler
}

// Role is the type representing a user's authorization level
type Role int

const (
	// Admin role, can do anything
	Admin Role = 0

	// Editor role, edit most objects
	Editor Role = 1

	// Reader role, readonly
	Reader Role = 2
)

// NotAuthorizedError represents an erorr in the authorization process
type NotAuthorizedError struct{}

func (e NotAuthorizedError) Error() string {
	return "user not authorized"
}

type userContextError struct{}

func (e userContextError) Error() string {
	return "user not found in context"
}
