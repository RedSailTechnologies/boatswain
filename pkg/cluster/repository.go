package cluster

import (
	"encoding/json"

	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
)

// Repository is the repository for dealing with cluster storage
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

// All gets all clusters, excluding deleted items
func (r *Repository) All() ([]*Cluster, error) {
	uuids, err := r.store.IDs(r.coll)
	if err != nil {
		return nil, err
	}

	clusters := make([]*Cluster, 0)
	for _, uuid := range uuids {
		storedEvents, err := r.store.GetEvents(r.coll, uuid)
		if err != nil {
			return nil, err
		}
		events, err := unmarshalEvents(storedEvents)
		if err != nil {
			return nil, err
		}
		cluster := Replay(events)
		if !cluster.destroyed {
			clusters = append(clusters, cluster)
		}
	}
	return clusters, nil
}

// Load reads out the cluster for the uuid given
func (r *Repository) Load(uuid string) (*Cluster, error) {
	events, err := r.store.GetEvents(r.coll, uuid)
	if err != nil {
		return nil, err
	}

	unmarshaled, err := unmarshalEvents(events)
	if err != nil {
		return nil, err
	}
	if len(unmarshaled) == 0 {
		return nil, ddd.NotFoundError{Entity: entityName}
	}

	cluster := Replay(unmarshaled)
	if cluster.destroyed {
		return nil, ddd.DestroyedError{Entity: entityName}
	}
	return cluster, nil
}

// Save persists the new events for the cluster given
func (r *Repository) Save(c *Cluster) error {
	version := r.store.GetVersion(r.coll, c.UUID())
	for i, ev := range c.Events()[version:] {
		v := i + version + 1
		d, err := json.Marshal(ev)
		if err != nil {
			return err
		}

		err = r.store.StoreEvent(r.coll, c.UUID(), ev.EventType(), string(d), v)
		if err != nil {
			return err
		}
	}
	return nil
}

func unmarshalEvents(events []*storage.StoredEvent) ([]ddd.Event, error) {
	unmarshaled := make([]ddd.Event, 0)
	var e interface{}
	for _, event := range events {
		switch event.Type {
		case entityName + "Created":
			e = &Created{}
		case entityName + "Destroyed":
			e = &Destroyed{}
		case entityName + "Updated":
			e = &Updated{}
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
