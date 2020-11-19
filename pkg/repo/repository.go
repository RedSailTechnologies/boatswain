package repo

import (
	"encoding/json"

	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
)

// Repository is the repository for dealing with repo storage
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

// All gets all repos, excluding deleted items
func (r *Repository) All() ([]*Repo, error) {
	uuids, err := r.store.IDs(r.coll)
	if err != nil {
		return nil, err
	}

	repos := make([]*Repo, 0)
	for _, uuid := range uuids {
		storedEvents, err := r.store.GetEvents(r.coll, uuid)
		if err != nil {
			return nil, err
		}
		events, err := unmarshalEvents(storedEvents)
		if err != nil {
			return nil, err
		}
		repo := Replay(events)
		if !repo.destroyed {
			repos = append(repos, repo)
		}
	}
	return repos, nil
}

// Load reads out the repo for the uuid given
func (r *Repository) Load(uuid string) (*Repo, error) {
	events, err := r.store.GetEvents(r.coll, uuid)
	if err != nil {
		return nil, err
	}

	unmarshaled, err := unmarshalEvents(events)
	if err != nil {
		return nil, err
	}
	if len(unmarshaled) == 0 {
		return nil, ddd.NotFoundError{Entity: "Repo"}
	}

	repo := Replay(unmarshaled)
	if repo.destroyed {
		return nil, ddd.DestroyedError{Entity: "Repo"}
	}
	return repo, nil
}

// Save persists the new events for the repo given
func (r *Repository) Save(repo *Repo) error {
	version := r.store.GetVersion(r.coll, repo.UUID())
	for i, ev := range repo.Events()[version:] {
		v := i + version + 1
		d, err := json.Marshal(ev)
		if err != nil {
			return err
		}

		err = r.store.StoreEvent(r.coll, repo.UUID(), ev.EventType(), string(d), v)
		if err != nil {
			return err
		}
	}
	return nil
}

func unmarshalEvents(events []*storage.StoredEvent) ([]ddd.Event, error) {
	unmarshaled := make([]ddd.Event, 0)
	for _, event := range events {
		switch event.Type {
		case "RepoCreated":
			var c Created
			err := json.Unmarshal([]byte(event.Data), &c)
			if err != nil {
				return nil, err
			}
			unmarshaled = append(unmarshaled, &c)
		case "RepoDestroyed":
			var d Destroyed
			err := json.Unmarshal([]byte(event.Data), &d)
			if err != nil {
				return nil, err
			}
			unmarshaled = append(unmarshaled, &d)
		case "RepoUpdated":
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
