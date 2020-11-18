package cluster

import "github.com/redsailtechnologies/boatswain/pkg/ddd"

// // Event is the interface for events
// type Event interface {
// 	isEvent()
// }

// Created is the event for when a new cluster is created
type Created struct {
	Timestamp int64
	UUID      string
	Name      string
	Endpoint  string
	Token     string
	Cert      string
}

// IsEvent marks this as an event
func (e Created) IsEvent() {}

// Destroyed is the event for when a cluster is destroyed
type Destroyed struct {
	Timestamp int64
}

// IsEvent marks this as an event
func (cu Destroyed) IsEvent() {}

// Updated is the event for when a cluster is updatedDestroyed
type Updated struct {
	Timestamp int64
	Name      string
	Endpoint  string
	Token     string
	Cert      string
}

// IsEvent marks this as an event
func (cu Updated) IsEvent() {}

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

// Replay recreates the aggregate from a series of events
func Replay(events []ddd.Event) *Cluster {
	c := &Cluster{}
	for _, event := range events {
		c.on(event)
	}
	return c
}

// Create handles create commands
func Create(uuid, name, endpoint, token, cert string, timestamp int64) (*Cluster, error) {
	if anyEmptyStrings(uuid, name, endpoint, token, cert) {
		return nil, ArgumentError
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
		return DestroyedError
	}
	c.on(&Destroyed{
		Timestamp: timestamp,
	})
	return nil
}

// Update handles update commands
func (c *Cluster) Update(name, endpoint, token, cert string, timestamp int64) error {
	if anyEmptyStrings(name, endpoint, token, cert) {
		return ArgumentError
	}
	if c.destroyed {
		return DestroyedError
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

// UUID returns this cluter's identifier
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
	return c.events
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

func anyEmptyStrings(strings ...string) bool {
	for _, str := range strings {
		if str == "" {
			return true
		}
	}
	return false
}

// ArgumentError represents an invalid argument passed to a command
var ArgumentError = argumentError{}

type argumentError struct{}

func (err argumentError) Error() string {
	return "all fields are required for a valid Cluster"
}

// DestroyedError represents an error when subsequent commands are called on a destroyed cluster
var DestroyedError = destroyedError{}

type destroyedError struct{}

func (err destroyedError) Error() string {
	return "Cluster cannot be modified further as it has been destroyed"
}
