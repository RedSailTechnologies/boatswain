package auth

import (
	"context"
	"net/http"
)

// Agent is the interface through which most of the auth package is utilized
type Agent interface {
	Authenticate(context.Context) (context.Context, error)
	Authorize(context.Context, Role) error
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

// AuthenticationError represents some error in the login process
type AuthenticationError struct{}

func (e AuthenticationError) Error() string {
	return "authentication error"
}

// NotAuthorizedError represents an erorr in the authorization process
type NotAuthorizedError struct{}

func (e NotAuthorizedError) Error() string {
	return "user not authorized"
}

type userContextError struct{}

func (e userContextError) Error() string {
	return "user not found in context"
}
