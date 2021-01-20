package cluster

import "github.com/redsailtechnologies/boatswain/pkg/ddd"

var entityName = "Cluster"

// Created is the event for when a new cluster is created
type Created struct {
	Timestamp int64
	UUID      string
	Name      string
	Endpoint  string
	Token     string
	Cert      string
}

// EventType marks this as an event
func (e Created) EventType() string {
	return entityName + "Created"
}

// Destroyed is the event for when a cluster is destroyed
type Destroyed struct {
	Timestamp int64
}

// EventType marks this as an event
func (e Destroyed) EventType() string {
	return entityName + "Destroyed"
}

// Updated is the event for when a cluster is updated
type Updated struct {
	Timestamp int64
	Name      string
	Endpoint  string
	Token     string
	Cert      string
}

// EventType marks this as an event
func (e Updated) EventType() string {
	return entityName + "Updated"
}

// Cluster represents a kubernetes cluster we are monitoring/deploying to
type Cluster struct {
	events    []ddd.Event
	version   int
	destroyed bool

	uuid     string
	name     string
	endpoint string
	token    string
	cert     string
}

// Replay recreates the cluster from a series of events
func Replay(events []ddd.Event) *Cluster {
	c := &Cluster{}
	for _, event := range events {
		c.on(event)
	}
	return c
}

// Create handles create commands
func Create(uuid, name, endpoint, token, cert string, timestamp int64) (*Cluster, error) {
	if uuid == "" {
		return nil, ddd.IDError{}
	}
	if err := validateFields(name, endpoint, token, cert); err != nil {
		return nil, err
	}

	c := &Cluster{}
	c.on(&Created{
		Timestamp: timestamp,
		UUID:      uuid,
		Name:      name,
		Endpoint:  endpoint,
		Token:     token,
		Cert:      cert,
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
func (c *Cluster) Update(name, endpoint, token, cert string, timestamp int64) error {
	if err := validateFields(name, endpoint, token, cert); err != nil {
		return err
	}
	if c.destroyed {
		return ddd.DestroyedError{Entity: entityName}
	}

	c.on(&Updated{
		Timestamp: timestamp,
		Name:      name,
		Endpoint:  endpoint,
		Token:     token,
		Cert:      cert,
	})
	return nil
}

// UUID returns this cluster's identifier
func (c *Cluster) UUID() string {
	return c.uuid
}

// Name returns this cluster's name
func (c *Cluster) Name() string {
	return c.name
}

// Endpoint returns this cluster's endpoint
func (c *Cluster) Endpoint() string {
	return c.endpoint
}

// Token returns this cluster's access token
func (c *Cluster) Token() string {
	return c.token
}

// Cert returns the cluster's cert info
func (c *Cluster) Cert() string {
	return c.cert
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
		c.endpoint = e.Endpoint
		c.token = e.Token
		c.cert = e.Cert
	case *Destroyed:
		c.destroyed = true
	case *Updated:
		c.name = e.Name
		c.endpoint = e.Endpoint
		c.token = e.Token
		c.cert = e.Cert
	}
}

func validateFields(name, endpoint, token, cert string) error {
	if name == "" {
		return ddd.RequiredArgumentError{Arg: "Name"}
	}
	if endpoint == "" {
		return ddd.RequiredArgumentError{Arg: "Endpoint"}
	}
	if token == "" {
		return ddd.RequiredArgumentError{Arg: "Token"}
	}
	if cert == "" {
		return ddd.RequiredArgumentError{Arg: "Cert"}
	}
	return nil
}
