package deployment

import "github.com/redsailtechnologies/boatswain/pkg/ddd"

var entityName = "Deployment"

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
func Replay(events []ddd.Event) ddd.Aggregate {
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

// Destroyed determines if this deployment has been destroyed
func (d *Deployment) Destroyed() bool {
	return d.destroyed
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
