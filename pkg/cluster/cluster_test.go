package cluster

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/stretchr/testify/assert"
)

func TestReplay(t *testing.T) {
	uuid := uuid.New().String()
	name := "acluster"
	endpoint := "http://cluster.cluster"
	token := "abc123"
	cert := "somecertdata"
	events := []ddd.Event{
		&Created{
			Timestamp: time.Now().Unix(),
			UUID:      uuid,
			Name:      "abc",
			Endpoint:  "something",
			Token:     "blah",
			Cert:      "certainly",
		},
		&Updated{
			Timestamp: time.Now().Unix(),
			Name:      name,
			Endpoint:  endpoint,
			Token:     token,
			Cert:      cert,
		},
	}

	sut := Replay(events)

	assert.Equal(t, uuid, sut.UUID())
	assert.Equal(t, name, sut.Name())
	assert.Equal(t, endpoint, sut.Endpoint())
	assert.Equal(t, token, sut.Token())
	assert.Equal(t, cert, sut.Cert())
	assert.Equal(t, 2, sut.Version())
	assert.Len(t, sut.Events(), 2)
}

func TestCreate(t *testing.T) {
	uuid := uuid.New().String()
	name := "acluster"
	endpoint := "http://cluster.cluster"
	token := "abc123"
	cert := "somecertdata"

	sut := Create(uuid, name, endpoint, token, cert, time.Now().Unix())

	assert.Equal(t, uuid, sut.UUID())
	assert.Equal(t, name, sut.Name())
	assert.Equal(t, endpoint, sut.Endpoint())
	assert.Equal(t, token, sut.Token())
	assert.Equal(t, cert, sut.Cert())
	assert.Equal(t, 1, sut.Version())
	assert.Len(t, sut.Events(), 1)
}

func TestDestroy(t *testing.T) {
	uuid := uuid.New().String()
	name := "acluster"
	endpoint := "http://cluster.cluster"
	token := "abc123"
	cert := "somecertdata"
	sut := Create(uuid, name, endpoint, token, cert, time.Now().Unix())

	sut.Destroy(time.Now().Unix())

	assert.Equal(t, true, sut.destroyed)
	assert.Len(t, sut.Events(), 2)
	assert.Equal(t, DestroyedError, sut.Destroy(time.Now().Unix()))
	assert.Equal(t, DestroyedError, sut.Update("", "", "", "", 0))
	assert.Len(t, sut.Events(), 2)
}

func TestUpdate(t *testing.T) {
	uuid := uuid.New().String()
	name := "acluster"
	endpoint := "http://cluster.cluster"
	token := "abc123"
	cert := "somecertdata"
	sut := Create(uuid, name, endpoint, token, cert, time.Now().Unix())

	name = "newname"
	endpoint = "http://new.cluster"
	token = "easy as"
	cert = "now with new cert data!"
	sut.Update(name, endpoint, token, cert, time.Now().Unix())

	assert.Equal(t, uuid, sut.UUID())
	assert.Equal(t, name, sut.Name())
	assert.Equal(t, endpoint, sut.Endpoint())
	assert.Equal(t, token, sut.Token())
	assert.Equal(t, cert, sut.Cert())
	assert.Equal(t, 2, sut.Version())
	assert.Len(t, sut.Events(), 2)
}

func TestDestroyedErrorMessage(t *testing.T) {
	assert.Equal(t, "Cluster cannot be modified further as it has been destroyed", DestroyedError.Error())
}
