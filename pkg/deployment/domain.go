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
func (d *Deployment) Destroy(timestamp int64) error {
	if d.destroyed {
		return ddd.DestroyedError{Entity: entityName}
	}
	d.on(&Destroyed{
		Timestamp: timestamp,
	})
	return nil
}

// Update handles update commands
func (d *Deployment) Update(name, repoID, branch, filePath string, timestamp int64) error {
	if err := validateFields(name, repoID, branch, filePath); err != nil {
		return err
	}
	if d.destroyed {
		return ddd.DestroyedError{Entity: entityName}
	}

	d.on(&Updated{
		Timestamp: timestamp,
		Name:      name,
		RepoID:    repoID,
		Branch:    branch,
		FilePath:  filePath,
	})
	return nil
}

// UUID returns this deployment's identifier
func (d *Deployment) UUID() string {
	return d.uuid
}

// Name returns this deployment's name
func (d *Deployment) Name() string {
	return d.name
}

// RepoID returns the uuid of the repo for this deployment's yaml file
func (d *Deployment) RepoID() string {
	return d.repoID
}

// Branch returns the branch of the repo for this deployment's yaml file
func (d *Deployment) Branch() string {
	return d.branch
}

// FilePath returns the file path of the repo for this deployment's yaml file
func (d *Deployment) FilePath() string {
	return d.filePath
}

// Events returns this deployment's event history
func (d *Deployment) Events() []ddd.Event {
	cp := make([]ddd.Event, len(d.events))
	copy(cp, d.events)
	return cp
}

// Version returns this deployment's version number (NOTE: aggregate version!)
func (d *Deployment) Version() int {
	return d.version
}

func (d *Deployment) on(event ddd.Event) {
	d.events = append(d.events, event)
	d.version++
	switch e := event.(type) {
	case *Created:
		d.uuid = e.UUID
		d.name = e.Name
		d.repoID = e.RepoID
		d.branch = e.Branch
		d.filePath = e.FilePath
	case *Destroyed:
		d.destroyed = true
	case *Updated:
		d.name = e.Name
		d.repoID = e.RepoID
		d.branch = e.Branch
		d.filePath = e.FilePath
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
