package cluster

import (
	"encoding/json"
	"errors"
	"fmt"

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

// All gets all clusters from the repo, excluding deleted items
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

	unmarshalled, err := unmarshalEvents(events)
	if err != nil {
		return nil, err
	}

	if len(unmarshalled) == 0 {
		return nil, NotFoundError
	}

	cluster := Replay(unmarshalled)
	if cluster.destroyed {
		return nil, errors.New("cluster has been destroyed")
	}
	return cluster, nil
}

// Save persists the cluster given
func (r *Repository) Save(c *Cluster) error {
	version := r.store.GetVersion(r.coll, c.UUID())
	for i, ev := range c.Events()[version:] {
		t, err := getEventType(ev)
		if err != nil {
			return err
		}
		d, err := json.Marshal(ev)
		if err != nil {
			return err
		}
		v := i + version + 1
		err = r.store.StoreEvent(r.coll, c.UUID(), t, string(d), v)
		if err != nil {
			return err
		}
	}
	return nil
}

// NotFoundError represents when a cluster cannot be found
var NotFoundError = notFoundError{}

type notFoundError struct{}

func (err notFoundError) Error() string {
	return "Cluster not found"
}

func getEventType(ev ddd.Event) (string, error) {
	switch fmt.Sprintf("%T", ev) {
	case "cluster.Created", "*cluster.Created":
		return "ClusterCreated", nil
	case "cluster.Destroyed", "*cluster.Destroyed":
		return "ClusterDestroyed", nil
	case "cluster.Updated", "*cluster.Updated":
		return "ClusterUpdated", nil
	default:
		return "", errors.New("event type not found")
	}
}

func unmarshalEvents(events []*storage.StoredEvent) ([]ddd.Event, error) {
	unmarshaled := make([]ddd.Event, 0)
	for _, event := range events {
		switch event.Type {
		case "ClusterCreated":
			var c Created
			err := json.Unmarshal([]byte(event.Data), &c)
			if err != nil {
				return nil, err
			}
			unmarshaled = append(unmarshaled, &c)
		case "ClusterDestroyed":
			var d Destroyed
			err := json.Unmarshal([]byte(event.Data), &d)
			if err != nil {
				return nil, err
			}
			unmarshaled = append(unmarshaled, &d)
		case "ClusterUpdated":
			var u Updated
			err := json.Unmarshal([]byte(event.Data), &u)
			if err != nil {
				return nil, err
			}
			unmarshaled = append(unmarshaled, &u)
		}
	}
	return unmarshaled, nil
}
