package storage

// Storage is the basic interface for storing of event streams
type Storage interface {
	All(coll string) ([]string, error)
	Load(coll, uuid string) ([]*StoredEvent, error)
	Save(coll, uuid, eventType, eventString string, version int) error
	Version(coll, uuid string) (int, error)
}

// StoredEvent represents an event in storage which includes its uuid, version, type, and event (as json)
type StoredEvent struct {
	UUID    string
	Version int
	Type    string
	Event   string
}
