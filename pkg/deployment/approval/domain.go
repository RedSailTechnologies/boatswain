package approval

import (
	"github.com/redsailtechnologies/boatswain/pkg/auth"
	"github.com/redsailtechnologies/boatswain/pkg/ddd"
)

var entityName = "Approval"

// Approval represents a manual signoff for a deployment
type Approval struct {
	events    []ddd.Event
	version   int
	destroyed bool

	uuid  string
	runID string

	completed  bool
	approved   bool
	overridden bool
	approver   *auth.User
	approvedAt int64

	allowedUsers []string
	allowedRoles []string
}

// Replay recreates the run from a series of events
func Replay(events []ddd.Event) ddd.Aggregate {
	a := &Approval{}
	for _, event := range events {
		a.on(event)
	}
	return a
}

// Create handles create commands for approvals
func Create(uuid, runID string, allowedUsers, allowedRoles []string) (*Approval, error) {
	if uuid == "" {
		return nil, ddd.IDError{}
	}
	if runID == "" {
		return nil, ddd.InvalidArgumentError{
			Arg: "runID",
			Val: "run uuid cannot be nil",
		}
	}
	if len(allowedUsers)+len(allowedRoles) == 0 {
		return nil, ApproverError{m: "at least one user or role approver required"}
	}

	a := &Approval{}
	a.on(&Created{
		UUID:         uuid,
		RunID:        runID,
		AllowedUsers: allowedUsers,
		AllowedRoles: allowedRoles,
	})
	return a, nil
}

// Complete handles approval/rejection commands for approvals
func (a *Approval) Complete(approve, override bool, user *auth.User, timestamp int64) error {
	valid := false
	for _, u := range a.allowedUsers {
		if user.Name == u {
			valid = true
		}
	}
	for _, r := range a.allowedRoles {
		if user.HasRole(r) {
			valid = true
		}
	}
	if override && user.HasRole(string(auth.Admin)) {
		valid = true
	}

	if !valid {
		return ApproverError{m: "user not authorized to approve this step"}
	}
	a.on(&Completed{
		Approved:   approve,
		Overridden: override,
		User:       user,
		Timestamp:  timestamp,
	})
	return nil
}

// Destroy handles destruction commands for when the approval is finished
func (a *Approval) Destroy() error {
	if a.destroyed {
		return ddd.DestroyedError{Entity: entityName}
	}
	a.on(&Destroyed{})
	return nil
}

// UUID gets the approval's uuid
func (a *Approval) UUID() string {
	return a.uuid
}

// Destroyed determines if this approval has been destroyed
func (a *Approval) Destroyed() bool {
	return a.destroyed // approvals can't be destroyed, but we still implement the interface
}

// RunUUID gets the approval's run id
func (a *Approval) RunUUID() string {
	return a.runID
}

// Completed determines if this approval has been approved or denied
func (a *Approval) Completed() bool {
	return a.completed
}

// Approved determines if this approval has been accepted
func (a *Approval) Approved() bool {
	return a.approved
}

// Overridden determines if this approval has been overriden
func (a *Approval) Overridden() bool {
	return a.overridden
}

// Approver gets the user who approved
func (a *Approval) Approver() (*auth.User, error) {
	if a.approver == nil {
		return nil, ApproverError{m: "not approved"}
	}
	return a.approver, nil
}

// Users gets users who can approve this approval
func (a *Approval) Users() []string {
	return a.allowedUsers
}

// Roles gets roles who can approve this approval
func (a *Approval) Roles() []string {
	return a.allowedRoles
}

// Events gets the run's event history
func (a *Approval) Events() []ddd.Event {
	cp := make([]ddd.Event, len(a.events))
	copy(cp, a.events)
	return cp
}

func (a *Approval) on(event ddd.Event) {
	a.events = append(a.events, event)
	a.version++
	switch e := event.(type) {
	case *Created:
		a.uuid = e.UUID
		a.runID = e.RunID
		a.allowedUsers = e.AllowedUsers
		a.allowedRoles = e.AllowedRoles
		a.approved = false
		a.overridden = false
	case *Completed:
		a.approver = e.User
		a.approvedAt = e.Timestamp
		a.approved = e.Approved
		a.overridden = e.Overridden
		a.completed = true
	case *Destroyed:
		a.destroyed = true
	}
}
