package cluster

import "github.com/redsailtechnologies/boatswain/pkg/ddd"

var entityName = "Cluster"

// Cluster represents a kubernetes cluster we are monitoring/deploying to
type Cluster struct {
	events    []ddd.Event
	version   int
	destroyed bool

	uuid  string
	name  string
	token string
}

// Replay recreates the cluster from a series of events
func Replay(events []ddd.Event) ddd.Aggregate {
	c := &Cluster{}
	for _, event := range events {
		c.on(event)
	}
	return c
}

// Create handles create commands
func Create(uuid, name, token string, timestamp int64) (*Cluster, error) {
	if uuid == "" {
		return nil, ddd.IDError{}
	}
	if name == "" {
		return nil, ddd.RequiredArgumentError{Arg: "Name"}
	}

	c := &Cluster{}
	c.on(&Created{
		Timestamp: timestamp,
		UUID:      uuid,
		Name:      name,
		Token:     token,
	})
	return c, nil
}

// Destroy handles destroy commands
func (c *Cluster) Destroy(timestamp int64) error {
	if c.destroyed {
		return ddd.DestroyedError{Entity: entityName}
	}
	c.on(&Destroyed{
		Timestamp: timestamp,
	})
	return nil
}

// Update handles update commands
func (c *Cluster) Update(name string, timestamp int64) error {
	if name == "" {
		return ddd.RequiredArgumentError{Arg: "Name"}
	}
	if c.destroyed {
		return ddd.DestroyedError{Entity: entityName}
	}

	c.on(&Updated{
		Timestamp: timestamp,
		Name:      name,
	})
	return nil
}

// UUID returns this cluster's identifier
func (c *Cluster) UUID() string {
	return c.uuid
}

// Destroyed determines if this deployment has been destroyed
func (c *Cluster) Destroyed() bool {
	return c.destroyed
}

// Name returns this cluster's name
func (c *Cluster) Name() string {
	return c.name
}

// Token returns this cluster's agent token
func (c *Cluster) Token() string {
	return c.token
}

// Events returns this cluster's event history
func (c *Cluster) Events() []ddd.Event {
	cp := make([]ddd.Event, len(c.events))
	copy(cp, c.events)
	return cp
}

// Version returns this cluster's version number (NOTE: aggregate version!)
func (c *Cluster) Version() int {
	return c.version
}

func (c *Cluster) on(event ddd.Event) {
	c.events = append(c.events, event)
	c.version++
	switch e := event.(type) {
	case *Created:
		c.uuid = e.UUID
		c.name = e.Name
		c.token = e.Token
	case *Destroyed:
		c.destroyed = true
	case *Updated:
		c.name = e.Name
	}
}
