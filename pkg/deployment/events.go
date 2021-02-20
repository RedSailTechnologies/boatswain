package deployment

import "github.com/redsailtechnologies/boatswain/pkg/ddd"

var eventTypes = map[string]func() ddd.Event{
	Created{}.EventType():   func() ddd.Event { return &Created{} },
	Destroyed{}.EventType(): func() ddd.Event { return &Destroyed{} },
	Updated{}.EventType():   func() ddd.Event { return &Updated{} },
}

// Created is the event for when a new deployment is created
type Created struct {
	Timestamp int64
	UUID      string
	Name      string
	Token     string
	RepoID    string
	Branch    string
	FilePath  string
}

// EventType marks this as an event
func (e Created) EventType() string {
	return entityName + "Created"
}

// Destroyed is the event for when a deployment is destroyed
type Destroyed struct {
	Timestamp int64
}

// EventType marks this as an event
func (e Destroyed) EventType() string {
	return entityName + "Destroyed"
}

// Updated is the event for when a deployment is updated
type Updated struct {
	Timestamp int64
	Name      string
	RepoID    string
	Branch    string
	FilePath  string
}

// EventType marks this as an event
func (e Updated) EventType() string {
	return entityName + "Updated"
}
