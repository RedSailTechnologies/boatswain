package cluster

import "github.com/redsailtechnologies/boatswain/pkg/ddd"

var eventTypes = map[string]ddd.Event{
	Created{}.EventType():   new(Created),
	Destroyed{}.EventType(): new(Destroyed),
	Updated{}.EventType():   new(Updated),
}

// Created is the event for when a new cluster is created
type Created struct {
	Timestamp int64
	UUID      string
	Name      string
	Endpoint  string
	Token     string
	Cert      string
}

// EventType marks this as an event
func (e Created) EventType() string {
	return entityName + "Created"
}

// Destroyed is the event for when a cluster is destroyed
type Destroyed struct {
	Timestamp int64
}

// EventType marks this as an event
func (e Destroyed) EventType() string {
	return entityName + "Destroyed"
}

// Updated is the event for when a cluster is updated
type Updated struct {
	Timestamp int64
	Name      string
	Endpoint  string
	Token     string
	Cert      string
}

// EventType marks this as an event
func (e Updated) EventType() string {
	return entityName + "Updated"
}
