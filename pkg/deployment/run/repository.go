package run

import (
	"strings"

	"github.com/redsailtechnologies/boatswain/pkg/storage"
)

// ReadRepository is the repository for dealing with run reads
type ReadRepository struct {
	r *storage.ReadRepository
}

// NewReadRepository creates a repository with the given storage
func NewReadRepository(s storage.Storage) *ReadRepository {
	return &ReadRepository{
		r: storage.NewReadRepository(strings.ToLower(entityName), eventTypes, Replay, s),
	}
}

// All gets all runs
func (rr *ReadRepository) All() ([]*Run, error) {
	results, err := rr.r.All()
	if err != nil {
		return nil, err
	}

	var runs []*Run
	for _, r := range results {
		runs = append(runs, r.(*Run))
	}
	return runs, nil
}

// Load gets one run
func (rr *ReadRepository) Load(uuid string) (*Run, error) {
	result, err := rr.r.Load(uuid)
	if err != nil {
		return nil, err
	}

	return result.(*Run), nil
}

type writeRepository struct {
	r *storage.WriteRepository
}

func newWriteRepository(s storage.Storage) *writeRepository {
	return &writeRepository{
		r: storage.NewWriteRepository(strings.ToLower(entityName), s),
	}
}

// Save persists new events for the run
func (wr *writeRepository) save(r *Run) error {
	return wr.r.Save(r)
}
