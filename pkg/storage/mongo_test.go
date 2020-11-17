package storage

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redsailtechnologies/boatswain/pkg/cluster"
	"github.com/stretchr/testify/assert"
)

func TestTemporary(t *testing.T) {
	conn := "mongodb://localhost:27017"
	mongo, err := NewMongo(conn, "tester")
	if err != nil {
		t.Error(err)
	}

	uuid := uuid.New().String()
	cluster := cluster.Create(uuid, "cluster", "endpoint", "token", "cert", time.Now().Unix())
	ev, _ := json.Marshal(cluster.Events()[0])
	err = mongo.Save("clusters", uuid, fmt.Sprintf("%T", cluster.Events()[0]), string(ev), 1)
	assert.Nil(t, err)
	err = mongo.Save("clusters", uuid, fmt.Sprintf("%T", cluster.Events()[0]), string(ev), 1)
	assert.NotNil(t, err)
	cluster.Update("something else", "endpoint", "token", "cert", time.Now().Unix())
	ev, _ = json.Marshal(cluster.Events()[1])
	err = mongo.Save("clusters", uuid, fmt.Sprintf("%T", cluster.Events()[0]), string(ev), 2)
	assert.Nil(t, err)
	events, err := mongo.Load("clusters", uuid)
	assert.NotNil(t, events)
	assert.Nil(t, err)
	all, err := mongo.All("clusters")
	assert.NotNil(t, all)
	version, err := mongo.Version("clusters", uuid)
	assert.Equal(t, 2, version)
	cluster.Update("something else again!", "endpoint", "token", "cert", time.Now().Unix())
	ev, _ = json.Marshal(cluster.Events()[2])
	err = mongo.Save("clusters", uuid, fmt.Sprintf("%T", cluster.Events()[0]), string(ev), 3)
	version, err = mongo.Version("clusters", uuid)
	assert.Equal(t, 3, version)
}
