package repo

import (
	"strings"

	"github.com/redsailtechnologies/boatswain/pkg/storage"
)

// ReadRepository is the repository for dealing with repo reads
type ReadRepository struct {
	r *storage.ReadRepository
}

// NewReadRepository creates a repository with the given storage
func NewReadRepository(s storage.Storage) *ReadRepository {
	return &ReadRepository{
		r: storage.NewReadRepository(strings.ToLower(entityName), eventTypes, Replay, s),
	}
}

// All gets all repos
func (rr *ReadRepository) All() ([]*Repo, error) {
	results, err := rr.r.All()
	if err != nil {
		return nil, err
	}

	var repos []*Repo
	for _, d := range results {
		repos = append(repos, d.(*Repo))
	}
	return repos, nil
}

// Load gets one repo
func (rr *ReadRepository) Load(uuid string) (*Repo, error) {
	result, err := rr.r.Load(uuid)
	if err != nil {
		return nil, err
	}

	return result.(*Repo), nil
}

type writeRepository struct {
	r *storage.WriteRepository
}

func newWriteRepository(s storage.Storage) *writeRepository {
	return &writeRepository{
		r: storage.NewWriteRepository(strings.ToLower(entityName), s),
	}
}

// Save persists new events for the repo
func (wr *writeRepository) save(d *Repo) error {
	return wr.r.Save(d)
}
