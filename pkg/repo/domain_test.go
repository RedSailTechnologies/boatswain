package repo

import (
	"testing"

	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/stretchr/testify/assert"
)

func TestEventTypes(t *testing.T) {
	assert.Equal(t, "RepoCreated", Created{}.EventType())
	assert.Equal(t, "RepoDestroyed", Destroyed{}.EventType())
	assert.Equal(t, "RepoUpdated", Updated{}.EventType())
}

func TestInvalidUUIDErrors(t *testing.T) {
	sut, err := Create("", "name", "endpoint", "", HELM, ddd.NewTimestamp())
	assert.Error(t, err)
	assert.Nil(t, sut)
}

func TestRequiredValidation(t *testing.T) {
	cases := map[string][]string{
		"Name":     []string{"", "endpoint"},
		"Endpoint": []string{"name", ""},
	}

	for k, v := range cases {
		id := ddd.NewUUID()
		name := v[0]
		endpoint := v[1]
		ti := ddd.NewTimestamp()
		sut, err := Create(id, name, endpoint, "", HELM, ti)
		assert.Error(t, err)
		assert.Equal(t, k, err.(ddd.RequiredArgumentError).Arg)
		assert.Nil(t, sut)
	}
}

func TestValidEndpoints(t *testing.T) {
	cases := []string{
		"http://repo.com",
		"https://repo.com",
	}

	for _, c := range cases {
		sut, err := Create(ddd.NewUUID(), "repo", c, "", HELM, ddd.NewTimestamp())
		assert.Nil(t, err)
		assert.NotNil(t, sut)
	}
}

func TestInvalidEndpoints(t *testing.T) {
	cases := []string{
		"starts.wrong",
		"short",
		"sevench",
	}

	for _, c := range cases {
		sut, err := Create(ddd.NewUUID(), "repo", c, "", HELM, ddd.NewTimestamp())
		assert.Error(t, err)
		assert.Equal(t, "Endpoint", err.(ddd.InvalidArgumentError).Arg)
		assert.Nil(t, sut)
	}
}

func TestReplay(t *testing.T) {
	uuid := ddd.NewUUID()
	name := "aRepo"
	endpoint := "https://repo.org"
	token := "abcdefg"
	events := []ddd.Event{
		&Created{
			Timestamp: ddd.NewTimestamp(),
			UUID:      uuid,
			Name:      "abc",
			Endpoint:  "http://repo.com",
			Token:     "",
		},
		&Updated{
			Timestamp: ddd.NewTimestamp(),
			Name:      name,
			Endpoint:  endpoint,
			Token:     token,
		},
	}

	sut := Replay(events).(*Repo)

	assert.Equal(t, uuid, sut.UUID())
	assert.Equal(t, name, sut.Name())
	assert.Equal(t, endpoint, sut.Endpoint())
	assert.Equal(t, token, sut.Token())
	assert.Equal(t, 2, sut.Version())
	assert.Len(t, sut.Events(), 2)
}

func TestCreate(t *testing.T) {
	uuid := ddd.NewUUID()
	name := "aRepo"
	endpoint := "http://repo.com"

	sut, err := Create(uuid, name, endpoint, "", HELM, ddd.NewTimestamp())

	assert.Nil(t, err)
	assert.Equal(t, uuid, sut.UUID())
	assert.Equal(t, name, sut.Name())
	assert.Equal(t, endpoint, sut.Endpoint())
	assert.Equal(t, 1, sut.Version())
	assert.Len(t, sut.Events(), 1)
}

func TestDestroy(t *testing.T) {
	uuid := ddd.NewUUID()
	name := "aRepo"
	endpoint := "http://cluster.cluster"

	sut, err := Create(uuid, name, endpoint, "", HELM, ddd.NewTimestamp())
	assert.Nil(t, err)

	err = sut.Destroy(ddd.NewTimestamp())
	assert.Nil(t, err)

	assert.Equal(t, true, sut.destroyed)
	assert.Len(t, sut.Events(), 2)
	assert.Equal(t, ddd.DestroyedError{Entity: "Repo"}, sut.Destroy(ddd.NewTimestamp()))
	assert.Equal(t, ddd.DestroyedError{Entity: "Repo"}, sut.Update("a", "http://a", "", HELM, 0))
	assert.Equal(t, ddd.RequiredArgumentError{Arg: "Endpoint"}, sut.Update("a", "", "", HELM, 0))
	assert.Len(t, sut.Events(), 2)
}

func TestUpdate(t *testing.T) {
	uuid := ddd.NewUUID()
	name := "aRepo"
	endpoint := "http://original.repo"
	token := ""

	sut, err := Create(uuid, name, endpoint, "", HELM, ddd.NewTimestamp())
	assert.Nil(t, err)

	name = "newname"
	endpoint = "https://new.repo"
	token = "anewtoken"
	sut.Update(name, endpoint, token, HELM, ddd.NewTimestamp())

	assert.Equal(t, uuid, sut.UUID())
	assert.Equal(t, name, sut.Name())
	assert.Equal(t, endpoint, sut.Endpoint())
	assert.Equal(t, token, sut.Token())
	assert.Equal(t, 2, sut.Version())
	assert.Len(t, sut.Events(), 2)
}
