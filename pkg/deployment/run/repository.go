package run

import (
	"encoding/json"

	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
)

// Repository is the repository for dealing with run storage
type Repository struct {
	coll  string
	store storage.Storage
}

// NewRepository creates a repository with the given storage
func NewRepository(coll string, store storage.Storage) *Repository {
	return &Repository{
		coll:  coll,
		store: store,
	}
}

// All gets all runs
// FIXME we need to implement a filter at the repo level here so we can
// get all runs by the deployment id, not every single one in the system
func (r *Repository) All() ([]*Run, error) {
	uuids, err := r.store.IDs(r.coll)
	if err != nil {
		return nil, err
	}

	runs := make([]*Run, 0)
	for _, uuid := range uuids {
		storedEvents, err := r.store.GetEvents(r.coll, uuid)
		if err != nil {
			return nil, err
		}
		events, err := unmarshalRunEvents(storedEvents)
		if err != nil {
			return nil, err
		}
		run := Replay(events)
		runs = append(runs, run)
	}
	return runs, nil
}

// Load reads out the run for the uuid given
func (r *Repository) Load(uuid string) (*Run, error) {
	events, err := r.store.GetEvents(r.coll, uuid)
	if err != nil {
		return nil, err
	}

	unmarshaled, err := unmarshalRunEvents(events)
	if err != nil {
		return nil, err
	}
	if len(unmarshaled) == 0 {
		return nil, ddd.NotFoundError{Entity: entityName}
	}

	return Replay(unmarshaled), nil
}

// Save persists the new events for the run given
func (r *Repository) Save(run *Run) error {
	version := r.store.GetVersion(r.coll, run.UUID())
	for i, ev := range run.Events()[version:] {
		v := i + version + 1
		d, err := json.Marshal(ev)
		if err != nil {
			return err
		}

		err = r.store.StoreEvent(r.coll, run.UUID(), ev.EventType(), string(d), v)
		if err != nil {
			return err
		}
	}
	return nil
}

func unmarshalRunEvents(events []*storage.StoredEvent) ([]ddd.Event, error) {
	unmarshaled := make([]ddd.Event, 0)
	var e interface{}
	for _, event := range events {
		switch event.Type {
		case entityName + "Created":
			e = &Created{}
		case entityName + "Started":
			e = &Started{}
		case entityName + "StepStarted":
			e = &StepStarted{}
		case entityName + "AppendLog":
			e = &AppendLog{}
		case entityName + "StepCompleted":
			e = &StepCompleted{}
		case entityName + "StepSkipped":
			e = &StepSkipped{}
		case entityName + "Completed":
			e = &Completed{}
		default:
			return nil, ddd.UnsupportedEventError{
				EventType: event.Type,
				Type:      entityName,
			}
		}
		err := json.Unmarshal([]byte(event.Data), &e)
		if err != nil {
			return nil, err
		}
		unmarshaled = append(unmarshaled, e.(ddd.Event))
	}
	return unmarshaled, nil
}
