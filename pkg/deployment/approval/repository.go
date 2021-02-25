package approval

import (
	"strings"

	"github.com/redsailtechnologies/boatswain/pkg/storage"
)

// ReadRepository is the repository for dealing with Approval reads
type ReadRepository struct {
	r *storage.ReadRepository
}

// NewReadRepository creates a repository with the given storage
func NewReadRepository(s storage.Storage) *ReadRepository {
	return &ReadRepository{
		r: storage.NewReadRepository(strings.ToLower(entityName), eventTypes, Replay, s),
	}
}

// All gets all Approvals
func (rr *ReadRepository) All() ([]*Approval, error) {
	results, err := rr.r.All()
	if err != nil {
		return nil, err
	}

	var Approvals []*Approval
	for _, a := range results {
		Approvals = append(Approvals, a.(*Approval))
	}
	return Approvals, nil
}

// Load gets one Approval
func (rr *ReadRepository) Load(uuid string) (*Approval, error) {
	result, err := rr.r.Load(uuid)
	if err != nil {
		return nil, err
	}

	return result.(*Approval), nil
}

// WriteRepository for creating/updating approvals
type WriteRepository struct {
	r *storage.WriteRepository
}

// NewWriteRepository builds the repository
func NewWriteRepository(s storage.Storage) *WriteRepository {
	return &WriteRepository{
		r: storage.NewWriteRepository(strings.ToLower(entityName), s),
	}
}

// Save persists new events for the Approval
func (wr *WriteRepository) Save(a *Approval) error {
	return wr.r.Save(a)
}
