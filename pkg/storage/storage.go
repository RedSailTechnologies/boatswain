package storage

// Storage is the basic interface for storing of event streams
type Storage interface {
	IDs(collection string) ([]string, error)
	GetEvents(collection, uuid string) ([]*StoredEvent, error)
	GetVersion(collection, uuid string) int
	StoreEvent(collection, uuid, eventType, eventData string, version int) error
}

// StoredEvent represents an event in storage which includes its uuid, version, type, and event (as json)
type StoredEvent struct {
	UUID    string
	Version int
	Type    string
	Data    string
}
