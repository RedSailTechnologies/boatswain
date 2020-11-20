package repo

import "github.com/redsailtechnologies/boatswain/pkg/ddd"

// Created is the event for when a new repo is created
type Created struct {
	Timestamp int64
	UUID      string
	Name      string
	Endpoint  string
}

// EventType marks this as an event
func (e Created) EventType() string {
	return entityName + "Created"
}

// Destroyed is the event for when a repo is destroyed
type Destroyed struct {
	Timestamp int64
}

// EventType marks this as an event
func (e Destroyed) EventType() string {
	return entityName + "Destroyed"
}

// Updated is the event for when a repo is updated
type Updated struct {
	Timestamp int64
	Name      string
	Endpoint  string
}

// EventType marks this as an event
func (e Updated) EventType() string {
	return entityName + "Updated"
}

var entityName = "Repo"

// Repo represents a repository, for now helm only
type Repo struct {
	events    []ddd.Event
	version   int
	destroyed bool

	uuid     string
	name     string
	endpoint string
}

// Replay recreates the repo from a series of events
func Replay(events []ddd.Event) *Repo {
	r := &Repo{}
	for _, event := range events {
		r.on(event)
	}
	return r
}

// Create handles create commands
func Create(uuid, name, endpoint string, timestamp int64) (*Repo, error) {
	if uuid == "" {
		return nil, ddd.IDError{}
	}
	err := validateFields(name, endpoint)
	if err != nil {
		return nil, err
	}

	r := &Repo{}
	r.on(&Created{
		Timestamp: timestamp,
		UUID:      uuid,
		Name:      name,
		Endpoint:  endpoint,
	})
	return r, nil
}

// Destroy handles destroy commands
func (r *Repo) Destroy(timestamp int64) error {
	if r.destroyed {
		return ddd.DestroyedError{Entity: entityName}
	}
	r.on(&Destroyed{
		Timestamp: timestamp,
	})
	return nil
}

// Update handles update commands
func (r *Repo) Update(name, endpoint string, timestamp int64) error {
	err := validateFields(name, endpoint)
	if err != nil {
		return err
	}
	if r.destroyed {
		return ddd.DestroyedError{Entity: entityName}
	}
	r.on(&Updated{
		Timestamp: timestamp,
		Name:      name,
		Endpoint:  endpoint,
	})
	return nil
}

// UUID gets the repo's unique id
func (r *Repo) UUID() string {
	return r.uuid
}

// Name gets the repo's name
func (r *Repo) Name() string {
	return r.name
}

// Endpoint gets the repo's endpoint
func (r *Repo) Endpoint() string {
	return r.endpoint
}

// Events gets all events from this repo
func (r *Repo) Events() []ddd.Event {
	return r.events
}

// Version gets all the
func (r *Repo) Version() int {
	return r.version
}

func (r *Repo) on(event ddd.Event) {
	r.events = append(r.events, event)
	r.version++
	switch e := event.(type) {
	case *Created:
		r.uuid = e.UUID
		r.name = e.Name
		r.endpoint = e.Endpoint
	case *Destroyed:
		r.destroyed = true
	case *Updated:
		r.name = e.Name
		r.endpoint = e.Endpoint
	}
}

func validateFields(name, endpoint string) error {
	if name == "" {
		return ddd.RequiredArgumentError{Arg: "Name"}
	}

	if endpoint == "" {
		return ddd.RequiredArgumentError{Arg: "Endpoint"}
	}

	if !validateRepoURL(endpoint) {
		return ddd.InvalidArgumentError{
			Arg: "Endpoint",
			Val: "must start with http:// or https://",
		}
	}

	return nil
}

func validateRepoURL(repoURL string) bool {
	if len(repoURL) < 7 {
		return false
	} else if repoURL[:7] == "http://" {
		return true
	}
	if len(repoURL) < 8 || repoURL[:8] != "https://" {
		return false
	}
	return true
}
