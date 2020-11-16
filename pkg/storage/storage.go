package storage

import "github.com/redsailtechnologies/boatswain/pkg/ddd"

// Storage is the basic interface for storing of event streams
type Storage interface {
	All() []string
	Load(uuid string) []ddd.Event
	Save(uuid string, version int, event ddd.Event) error
}
