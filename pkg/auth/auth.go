package auth

import (
	"context"
	"net/http"
)

// Agent is the interface through which most of the auth package is utilized
type Agent interface {
	Authenticate(context.Context) (context.Context, error)
	Authorize(context.Context, Role) error
	User(context.Context) User
	Wrap(http.Handler) http.Handler
}

// Role is the type representing a user's authorization level
type Role string

const (
	// Admin role, can do anything
	Admin Role = "Admin"

	// Editor role, edit most objects
	Editor Role = "Editor"

	// Reader role, readonly
	Reader Role = "Reader"
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
