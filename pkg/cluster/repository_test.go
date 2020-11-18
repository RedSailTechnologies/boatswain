package cluster

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestTemporary(t *testing.T) {
	conn := "mongodb://localhost:27017"
	mongo, err := storage.NewMongo(conn, "tester")
	if err != nil {
		return
	}
	repo := NewRepository("clusters", mongo)
	repo.Load("doesn'texist")
	clusters, err := repo.All()
	if err != nil {
		return
	}
	assert.NotNil(t, clusters)
	c, _ := Create(uuid.New().String(), "another cluster", "some endpoint", "a token", "a cert", time.Now().Unix())
	err = repo.Save(c)
	assert.Nil(t, err)
	d, err := repo.Load(c.UUID())
	assert.NotNil(t, d)
	clusters, err = repo.All()
	if err != nil {
		return
	}
}
