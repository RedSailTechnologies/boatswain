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
	token := ddd.NewUUID()
	ti := ddd.NewTimestamp()
	sut, err := Create(id, name, ddd.NewUUID(), ti)
	assert.Error(t, err)
	assert.Equal(t, k, err.(ddd.RequiredArgumentError).Arg)
	assert.Nil(t, sut)
	cases := map[string][]string{
		"Name": []string{""},
	}

	for k, v := range cases {
	}
}

func TestReplay(t *testing.T) {
	uuid := ddd.NewUUID()
	name := "acluster"
	endpoint := "http://cluster.cluster"
	token := "abc123"
	cert := "somecertdata"
	events := []ddd.Event{
		&Created{
			Timestamp: ddd.NewTimestamp(),
			UUID:      uuid,
			Name:      "abc",
			Endpoint:  "something",
			Token:     "blah",
			Cert:      "certainly",
		},
		&Updated{
			Timestamp: ddd.NewTimestamp(),
			Name:      name,
			Endpoint:  endpoint,
			Token:     token,
			Cert:      cert,
		},
	}

	sut := Replay(events).(*Cluster)

	assert.Equal(t, uuid, sut.UUID())
	assert.Equal(t, name, sut.Name())
	assert.Equal(t, endpoint, sut.Endpoint())
	assert.Equal(t, token, sut.Token())
	assert.Equal(t, cert, sut.Cert())
	assert.Equal(t, 2, sut.Version())
	assert.Len(t, sut.Events(), 2)
}

func TestCreate(t *testing.T) {
	uuid := ddd.NewUUID()
	name := "acluster"
	endpoint := "http://cluster.cluster"
	token := "abc123"
	cert := "somecertdata"

	sut, err := Create(uuid, name, endpoint, token, cert, ddd.NewTimestamp())

	assert.Nil(t, err)
	assert.Equal(t, uuid, sut.UUID())
	assert.Equal(t, name, sut.Name())
	assert.Equal(t, endpoint, sut.Endpoint())
	assert.Equal(t, token, sut.Token())
	assert.Equal(t, cert, sut.Cert())
	assert.Equal(t, 1, sut.Version())
	assert.Len(t, sut.Events(), 1)
}

func TestDestroy(t *testing.T) {
	uuid := ddd.NewUUID()
	name := "acluster"
	endpoint := "http://cluster.cluster"
	token := "abc123"
	cert := "somecertdata"

	sut, err := Create(uuid, name, endpoint, token, cert, ddd.NewTimestamp())
	assert.Nil(t, err)

	err = sut.Destroy(ddd.NewTimestamp())
	assert.Nil(t, err)

	assert.Equal(t, true, sut.destroyed)
	assert.Len(t, sut.Events(), 2)
	assert.Equal(t, ddd.DestroyedError{Entity: "Cluster"}, sut.Destroy(ddd.NewTimestamp()))
	assert.Equal(t, ddd.DestroyedError{Entity: "Cluster"}, sut.Update("a", "b", "c", "d", 0))
	assert.Equal(t, ddd.RequiredArgumentError{Arg: "Endpoint"}, sut.Update("a", "", "b", "c", 0))
	assert.Len(t, sut.Events(), 2)
}

func TestUpdate(t *testing.T) {
	uuid := ddd.NewUUID()
	name := "acluster"
	endpoint := "http://cluster.cluster"
	token := "abc123"
	cert := "somecertdata"

	sut, err := Create(uuid, name, endpoint, token, cert, ddd.NewTimestamp())
	assert.Nil(t, err)

	name = "newname"
	endpoint = "http://new.cluster"
	token = "easy as"
	cert = "now with new cert data!"
	sut.Update(name, endpoint, token, cert, ddd.NewTimestamp())

	assert.Equal(t, uuid, sut.UUID())
	assert.Equal(t, name, sut.Name())
	assert.Equal(t, endpoint, sut.Endpoint())
	assert.Equal(t, token, sut.Token())
	assert.Equal(t, cert, sut.Cert())
	assert.Equal(t, 2, sut.Version())
	assert.Len(t, sut.Events(), 2)
}
