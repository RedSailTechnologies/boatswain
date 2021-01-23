package storage

import (
	"encoding/json"
	"strings"

	"github.com/redsailtechnologies/boatswain/pkg/ddd"
)

// ReadRepository is a generic repository for reads
type ReadRepository struct {
	name       string
	eventTypes map[string]func() ddd.Event
	replay     func([]ddd.Event) ddd.Aggregate
	store      Storage
}

// NewReadRepository gets a new generic repository for reads
func NewReadRepository(n string,
	et map[string]func() ddd.Event,
	r func([]ddd.Event) ddd.Aggregate,
	s Storage) *ReadRepository {
	return &ReadRepository{
		name:       n,
		eventTypes: et,
		replay:     r,
		store:      s,
	}
}

// All gets all aggregates
func (r *ReadRepository) All() ([]ddd.Aggregate, error) {
	uuids, err := r.store.IDs(r.name)
	if err != nil {
		return nil, err
	}

	var aggregates []ddd.Aggregate
	for _, uuid := range uuids {
		storedEvents, err := r.store.GetEvents(r.name, uuid)
		if err != nil {
			return nil, err
		}
		events, err := r.unmarshal(storedEvents)
		if err != nil {
			return nil, err
		}
		aggregate := r.replay(events)
		if !aggregate.Destroyed() {
			aggregates = append(aggregates, aggregate)
		}
	}
	return aggregates, nil
}

// Load reads out a specific aggregate with the given uuid
func (r *ReadRepository) Load(uuid string) (ddd.Aggregate, error) {
	events, err := r.store.GetEvents(r.name, uuid)
	if err != nil {
		return nil, err
	}

	unmarshaled, err := r.unmarshal(events)
	if err != nil {
		return nil, err
	}
	if len(unmarshaled) == 0 {
		return nil, ddd.NotFoundError{Entity: strings.Title(r.name)}
	}

	aggregate := r.replay(unmarshaled)
	if aggregate.Destroyed() {
		return nil, ddd.DestroyedError{Entity: strings.Title(r.name)}
	}
	return aggregate, nil
}

func (r *ReadRepository) unmarshal(events []*StoredEvent) ([]ddd.Event, error) {
	unmarshaled := make([]ddd.Event, 0)
	for _, event := range events {
		if e, ok := r.eventTypes[event.Type]; !ok {
			return nil, ddd.UnsupportedEventError{
				EventType: event.Type,
				Type:      strings.Title(r.name),
			}
		} else {
			out := e()
			err := json.Unmarshal([]byte(event.Data), &out)
			if err != nil {
				return nil, err
			}
			unmarshaled = append(unmarshaled, out.(ddd.Event))
		}

	}
	return unmarshaled, nil
}

// WriteRepository is a generic repository for reads/writes
type WriteRepository struct {
	name  string
	store Storage
}

// NewWriteRepository gets a new generic repository for reads/writes
func NewWriteRepository(n string, s Storage) *WriteRepository {
	return &WriteRepository{
		name:  n,
		store: s,
	}
}

// Save persists the new events for the aggregate
func (r *WriteRepository) Save(a ddd.Aggregate) error {
	version := r.store.GetVersion(r.name, a.UUID())
	for i, ev := range a.Events()[version:] {
		v := i + version + 1
		e, err := json.Marshal(ev)
		if err != nil {
			return err
		}

		err = r.store.StoreEvent(r.name, a.UUID(), ev.EventType(), string(e), v)
		if err != nil {
			return err
		}
	}
	return nil
}
