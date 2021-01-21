package cluster

import (
	"strings"

	"github.com/redsailtechnologies/boatswain/pkg/storage"
)

// ReadRepository is the repository for dealing with cluster reads
type ReadRepository struct {
	r *storage.ReadRepository
}

// NewReadRepository creates a repository with the given storage
func NewReadRepository(s storage.Storage) *ReadRepository {
	return &ReadRepository{
		r: storage.NewReadRepository(strings.ToLower(entityName), eventTypes, Replay, s),
	}
}

// All gets all clusters
func (rr *ReadRepository) All() ([]*Cluster, error) {
	results, err := rr.r.All()
	if err != nil {
		return nil, err
	}

	var clusters []*Cluster
	for _, d := range results {
		clusters = append(clusters, d.(*Cluster))
	}
	return clusters, nil
}

// Load gets one cluster
func (rr *ReadRepository) Load(uuid string) (*Cluster, error) {
	result, err := rr.r.Load(uuid)
	if err != nil {
		return nil, err
	}

	return result.(*Cluster), nil
}

type writeRepository struct {
	r *storage.WriteRepository
}

func newWriteRepository(s storage.Storage) *writeRepository {
	return &writeRepository{
		r: storage.NewWriteRepository(strings.ToLower(entityName), s),
	}
}

// Save persists new events for the cluster
func (wr *writeRepository) save(d *Cluster) error {
	return wr.r.Save(d)
}
