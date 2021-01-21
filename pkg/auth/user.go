package auth

import (
	"fmt"
)

// User represents our representation of an oidc user
type User struct {
	token string

	Name         string   `json:"name"`
	Email        string   `json:"email"`
	Subject      string   `json:"sub"`
	NotBefore    int64    `json:"nbf"`
	NotOnOrAfter int64    `json:"exp"`
	IssuedAt     int64    `json:"iat"`
	Scope        string   `json:"scp"`
	Roles        []string `json:"roles"`
}

func (u *User) hasRole(r string) bool {
	for _, role := range u.Roles {
		if role == r {
			return true
		}
	}
	return false
}

func (u *User) validateScope(s string) error {
	if u.Scope != s {
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
