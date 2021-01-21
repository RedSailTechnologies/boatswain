package repo

import "github.com/redsailtechnologies/boatswain/pkg/ddd"

var entityName = "Repo"

// Repo represents a repository, for now helm only
type Repo struct {
	events    []ddd.Event
	version   int
	destroyed bool

	uuid     string
	name     string
	endpoint string
	repoType Type
}

// Replay recreates the repo from a series of events
func Replay(events []ddd.Event) ddd.Aggregate {
	r := &Repo{}
	for _, event := range events {
		r.on(event)
	}
	return r
}

// Create handles create commands
func Create(uuid, name, endpoint string, t Type, timestamp int64) (*Repo, error) {
	if uuid == "" {
		return nil, ddd.IDError{}
	}
	if err := validateFields(name, endpoint); err != nil {
		return nil, err
	}

	r := &Repo{}
	r.on(&Created{
		Timestamp: timestamp,
		UUID:      uuid,
		Name:      name,
		Endpoint:  endpoint,
		Type:      t,
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
func (r *Repo) Update(name, endpoint string, t Type, timestamp int64) error {
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
		Type:      t,
	})
	return nil
}

// UUID gets the repo's unique id
func (r *Repo) UUID() string {
	return r.uuid
}

// Destroyed determines if this repo has been destroyed
func (r *Repo) Destroyed() bool {
	return r.destroyed
}

// Name gets the repo's name
func (r *Repo) Name() string {
	return r.name
}

// Endpoint gets the repo's endpoint
func (r *Repo) Endpoint() string {
	return r.endpoint
}

// Type gets the repo type (helm, git, etc.)
func (r *Repo) Type() Type {
	return r.repoType
}

// Events gets all events from this repo
func (r *Repo) Events() []ddd.Event {
	cp := make([]ddd.Event, len(r.events))
	copy(cp, r.events)
	return cp
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
		r.repoType = e.Type
	case *Destroyed:
		r.destroyed = true
	case *Updated:
		r.name = e.Name
		r.endpoint = e.Endpoint
		r.repoType = e.Type
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
