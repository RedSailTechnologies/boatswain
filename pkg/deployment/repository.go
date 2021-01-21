package deployment

import (
	"strings"

	"github.com/redsailtechnologies/boatswain/pkg/storage"
)

// ReadRepository is the repository for dealing with deployment reads
type ReadRepository struct {
	r *storage.ReadRepository
}

// NewReadRepository creates a repository with the given storage
func NewReadRepository(s storage.Storage) *ReadRepository {
	return &ReadRepository{
		r: storage.NewReadRepository(strings.ToLower(entityName), eventTypes, Replay, s),
	}
}

// All gets all deployments
func (rr *ReadRepository) All() ([]*Deployment, error) {
	results, err := rr.r.All()
	if err != nil {
		return nil, err
	}

	var deployments []*Deployment
	for _, d := range results {
		deployments = append(deployments, d.(*Deployment))
	}
	return deployments, nil
}

// Load gets one deployment
func (rr *ReadRepository) Load(uuid string) (*Deployment, error) {
	result, err := rr.r.Load(uuid)
	if err != nil {
		return nil, err
	}

	return result.(*Deployment), nil
}

type writeRepository struct {
	r *storage.WriteRepository
}

func newWriteRepository(s storage.Storage) *writeRepository {
	return &writeRepository{
		r: storage.NewWriteRepository(strings.ToLower(entityName), s),
	}
}

// Save persists new events for the deployment
func (wr *writeRepository) save(d *Deployment) error {
	return wr.r.Save(d)
}
