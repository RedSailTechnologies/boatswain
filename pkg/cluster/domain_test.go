package cluster

import (
	"testing"

	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/stretchr/testify/assert"
)

func TestEventTypes(t *testing.T) {
	assert.Equal(t, "ClusterCreated", Created{}.EventType())
	assert.Equal(t, "ClusterDestroyed", Destroyed{}.EventType())
	assert.Equal(t, "ClusterUpdated", Updated{}.EventType())
}

func TestInvalidUUIDErrors(t *testing.T) {
	sut, err := Create("", "name", "token", ddd.NewTimestamp())
	assert.Error(t, err)
	assert.Nil(t, sut)
}

func TestValidation(t *testing.T) {
	id := ddd.NewUUID()
	name := ""
	ti := ddd.NewTimestamp()
	sut, err := Create(id, name, ddd.NewUUID(), ti)
	assert.Error(t, err)
	assert.Equal(t, "name", err.(ddd.RequiredArgumentError).Arg)
	assert.Nil(t, sut)
}

func TestReplay(t *testing.T) {
	uuid := ddd.NewUUID()
	name := "acluster"
	token := "abc123"
	events := []ddd.Event{
		&Created{
			Timestamp: ddd.NewTimestamp(),
			UUID:      uuid,
			Name:      "abc",
			Token:     token,
		},
		&Updated{
			Timestamp: ddd.NewTimestamp(),
			Name:      name,
		},
	}

	sut := Replay(events).(*Cluster)

	assert.Equal(t, uuid, sut.UUID())
	assert.Equal(t, name, sut.Name())
	assert.Equal(t, token, sut.Token())
	assert.Equal(t, 2, sut.Version())
	assert.Len(t, sut.Events(), 2)
}

func TestCreate(t *testing.T) {
	uuid := ddd.NewUUID()
	name := "acluster"
	token := "abc123"

	sut, err := Create(uuid, name, token, ddd.NewTimestamp())

	assert.Nil(t, err)
	assert.Equal(t, uuid, sut.UUID())
	assert.Equal(t, name, sut.Name())
	assert.Equal(t, token, sut.Token())
	assert.Equal(t, 1, sut.Version())
	assert.Len(t, sut.Events(), 1)
}

func TestDestroy(t *testing.T) {
	uuid := ddd.NewUUID()
	name := "acluster"
	token := "abc123"

	sut, err := Create(uuid, name, token, ddd.NewTimestamp())
	assert.Nil(t, err)

	err = sut.Destroy(ddd.NewTimestamp())
	assert.Nil(t, err)

	assert.Equal(t, true, sut.destroyed)
	assert.Len(t, sut.Events(), 2)
	assert.Equal(t, ddd.DestroyedError{Entity: "Cluster"}, sut.Destroy(ddd.NewTimestamp()))
	assert.Equal(t, ddd.DestroyedError{Entity: "Cluster"}, sut.Update("a", 0))
	assert.Len(t, sut.Events(), 2)
}

func TestUpdate(t *testing.T) {
	uuid := ddd.NewUUID()
	name := "acluster"
	token := "abc123"

	sut, err := Create(uuid, name, token, ddd.NewTimestamp())
	assert.Nil(t, err)

	name = "newname"
	sut.Update(name, ddd.NewTimestamp())

	assert.Equal(t, uuid, sut.UUID())
	assert.Equal(t, name, sut.Name())
	assert.Equal(t, token, sut.Token())
	assert.Equal(t, 2, sut.Version())
	assert.Len(t, sut.Events(), 2)
}
