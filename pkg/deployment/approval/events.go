package approval

import (
	"github.com/redsailtechnologies/boatswain/pkg/auth"
	"github.com/redsailtechnologies/boatswain/pkg/ddd"
)

var eventTypes = map[string]func() ddd.Event{
	Created{}.EventType():   func() ddd.Event { return &Created{} },
	Completed{}.EventType(): func() ddd.Event { return &Completed{} },
	Destroyed{}.EventType(): func() ddd.Event { return &Destroyed{} },
}

// Created is the event for when an approval is created
type Created struct {
	UUID         string
	RunID        string
	AllowedUsers []string
	AllowedRoles []string
}

// EventType marks this as an event
func (e Created) EventType() string {
	return entityName + "Created"
}

// Completed is the event for when an approval is accepted
type Completed struct {
	Approved   bool
	Overridden bool
	User       *auth.User
	Timestamp  int64
}

// EventType marks this as an event
func (e Completed) EventType() string {
	return entityName + "Approved"
}

// Destroyed is the event for when an approval is finalized
type Destroyed struct{}

// EventType marks this as an event
func (e Destroyed) EventType() string {
	return entityName + "Destroyed"
}
