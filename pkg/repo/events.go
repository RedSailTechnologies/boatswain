package repo

// Created is the event for when a new repo is created
type Created struct {
	Timestamp int64
	UUID      string
	Name      string
	Endpoint  string
	Type      Type
}

// EventType marks this as an event
func (e Created) EventType() string {
	return entityName + "Created"
}

// Destroyed is the event for when a repo is destroyed
type Destroyed struct {
	Timestamp int64
}

// EventType marks this as an event
func (e Destroyed) EventType() string {
	return entityName + "Destroyed"
}

// Updated is the event for when a repo is updated
type Updated struct {
	Timestamp int64
	Name      string
	Endpoint  string
	Type      Type
}

// EventType marks this as an event
func (e Updated) EventType() string {
	return entityName + "Updated"
}
