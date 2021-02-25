package auth

import (
	"fmt"
)

// User is the external (to the system) representation of a user
type User struct {
	Name    string
	Email   string
	Subject string
	Roles   []Role
}

// HasRole checks if a user has a role represented by the string passed
func (u *User) HasRole(r string) bool {
	role := Role(r)
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

type user struct {
	token string

	Name         string   `json:"name"`
	Email        string   `json:"email"`
	Subject      string   `json:"sub"`
	IssuedAt     int64    `json:"iat"`
	NotBefore    int64    `json:"nbf"`
	NotOnOrAfter int64    `json:"exp"`
	Scope        string   `json:"scp"`
	Roles        []string `json:"roles"`
}

func (u *user) hasRole(r string) bool {
	for _, role := range u.Roles {
		if role == r {
			return true
		}
	}
	return false
}

func (u *user) validateScope(s string) error {
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
