package auth

import (
	"context"
	"errors"
	"net/http"
)

// Agent is the interface through which most of the auth package is utilized
type Agent interface {
	Authenticate(context.Context) (context.Context, error)
	Authorize(context.Context, Role) error
	User(context.Context) User
	Roles(User) []Role
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

// ToRole takes a string and tries to convert it to a role
func ToRole(r string) (Role, error) {
	switch r {
	case "Admin":
		return Admin, nil
	case "Editor":
		return Editor, nil
	case "Reader":
		return Reader, nil
	default:
		return -1, errors.New("role not found")
	}
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
