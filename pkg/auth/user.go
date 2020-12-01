package auth

import (
	"context"
	"fmt"
)

var userKey = new(int)

// User represents our representation of an oidc user
type User struct {
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	Subject      string   `json:"sub"`
	NotBefore    int64    `json:"nbf"`
	NotOnOrAfter int64    `json:"exp"`
	IssuedAt     int64    `json:"iat"`
	Scope        string   `json:"scp"`
	Scopes       []string `json:"scope"`
	Roles        []string `json:"roles"`
}

// AddToContext adds the user to the context
func (u *User) AddToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, userKey, u)
}

// UserFromContext gets the user from the context
func UserFromContext(ctx context.Context) *User {
	userVal := ctx.Value(userKey)
	if userVal == nil {
		return nil
	}
	fmt.Printf("%T\n", userVal)
	user := userVal.(*User)
	return user
}

// IsAdmin checks to see if the user is in the admin role
func (u *User) IsAdmin() bool {
	for _, role := range u.Roles {
		if role == config.adminRole {
			return true
		}
	}
	return false
}

// IsEditor checks to see if the user is in the editor role
func (u *User) IsEditor() bool {
	for _, role := range u.Roles {
		if role == config.editorRole {
			return true
		}
	}
	return false
}

// IsReader checks to see if the user is in the reader role
func (u *User) IsReader() bool {
	for _, role := range u.Roles {
		if role == config.readerRole {
			return true
		}
	}
	return false
}

// ValidateScope checks the user scope for this app
func (u *User) ValidateScope() error {
	if u.Scope != config.scope {
		for _, scope := range u.Scopes {
			if scope == config.scope {
				return nil
			}
		}
		return &userScopeError{scope: u.Scope}
	}
	return nil
}

type userScopeError struct {
	scope string
}

func (err userScopeError) Error() string {
	return fmt.Sprintf("user does not have the %s scope", err.scope)
}
