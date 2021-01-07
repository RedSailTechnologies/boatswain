package deployment

import "github.com/redsailtechnologies/boatswain/pkg/ddd"

var entityName = "Deployment"

// Created is the event for when a new deployment is created
type Created struct {
	Timestamp int64
	UUID      string
	Name      string
	RepoID    string
	Branch    string
	FilePath  string
}

// EventType marks this as an event
func (e Created) EventType() string {
	return entityName + "Created"
}

// Destroyed is the event for when a deployment is destroyed
type Destroyed struct {
	Timestamp int64
}

// EventType marks this as an event
func (e Destroyed) EventType() string {
	return entityName + "Destroyed"
}

// Updated is the event for when a deployment is updated
type Updated struct {
	Timestamp int64
	Name      string
	RepoID    string
	Branch    string
	FilePath  string
}

// EventType marks this as an event
func (e Updated) EventType() string {
	return entityName + "Updated"
}

// Deployment represents the basic information about a deployment
type Deployment struct {
	events    []ddd.Event
	version   int
	destroyed bool

	uuid     string
	name     string
	repoID   string
	branch   string
	filePath string
}

// Replay recreates the deployment from a series of events
func Replay(events []ddd.Event) *Deployment {
	d := &Deployment{}
	for _, event := range events {
		d.on(event)
	}
	return d
}

// Create handles create commands
func Create(uuid, name, repoID, branch, filePath string, timestamp int64) (*Deployment, error) {
	if uuid == "" {
		return nil, ddd.IDError{}
	}
	if err := validateFields(name, repoID, branch, filePath); err != nil {
		return nil, err
	}

	d := &Deployment{}
	d.on(&Created{
		Timestamp: timestamp,
		UUID:      uuid,
		Name:      name,
		RepoID:    repoID,
		Branch:    branch,
		FilePath:  filePath,
	})
	return d, nil
}

// Destroy handles destroy commands
func (c *Deployment) Destroy(timestamp int64) error {
	if c.destroyed {
		return ddd.DestroyedError{Entity: entityName}
	}
	c.on(&Destroyed{
		Timestamp: timestamp,
	})
	return nil
}

// Update handles update commands
func (c *Deployment) Update(name, repoID, branch, filePath string, timestamp int64) error {
	if err := validateFields(name, repoID, branch, filePath); err != nil {
		return err
	}
	if c.destroyed {
		return ddd.DestroyedError{Entity: entityName}
	}

	c.on(&Updated{
		Timestamp: timestamp,
		Name:      name,
		RepoID:    repoID,
		Branch:    branch,
		FilePath:  filePath,
	})
	return nil
}

// UUID returns this deployment's identifier
func (c *Deployment) UUID() string {
	return c.uuid
}

// Name returns this deployment's name
func (c *Deployment) Name() string {
	return c.name
}

// RepoID returns the uuid of the repo for this deployment's yaml file
func (c *Deployment) RepoID() string {
	return c.repoID
}

// Branch returns the branch of the repo for this deployment's yaml file
func (c *Deployment) Branch() string {
	return c.branch
}

// FilePath returns the file path of the repo for this deployment's yaml file
func (c *Deployment) FilePath() string {
	return c.name
}

// Events returns this deployment's event history
func (c *Deployment) Events() []ddd.Event {
	return c.events
}

// Version returns this deployment's version number (NOTE: aggregate version!)
func (c *Deployment) Version() int {
	return c.version
}

func (c *Deployment) on(event ddd.Event) {
	c.events = append(c.events, event)
	c.version++
	switch e := event.(type) {
	case *Created:
		c.uuid = e.UUID
		c.name = e.Name
		c.repoID = e.RepoID
		c.branch = e.Branch
		c.filePath = e.FilePath
	case *Destroyed:
		c.destroyed = true
	case *Updated:
		c.name = e.Name
		c.repoID = e.RepoID
		c.branch = e.Branch
		c.filePath = e.FilePath
	}
}

func validateFields(name, repoID, branch, filePath string) error {
	if name == "" {
		return ddd.RequiredArgumentError{Arg: "Name"}
	}
	if repoID == "" {
		return ddd.RequiredArgumentError{Arg: "RepoID"}
	}
	if branch == "" {
		return ddd.RequiredArgumentError{Arg: "Branch"}
	}
	if filePath == "" {
		return ddd.RequiredArgumentError{Arg: "FilePath"}
	}
	return nil
}
