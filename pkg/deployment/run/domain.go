package run

import (
	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/deployment/template"
)

var entityName = "Run"

/* Events:
Created
StepStarted
StepCompleted
ActionStarted
ActionCompleted
*/

// Created is the event for when a new run is created
type Created struct {
	Timestamp      int64
	UUID           string
	DeploymentUUID string
}

// EventType marks this as an event
func (e Created) EventType() string {
	return entityName + "Created"
}

// Run represents a single execution of a deployment
type Run struct {
	events  []ddd.Event
	version int

	uuid           string
	deploymentUUID string
	status         status

	template       template.Deployment
	templateResult error
	verifyResult   error

	trigger       todo.Trigger
	triggerResult error

	stepResults []stepResult
}

type stepResult struct {
	actionResults []actionResult
}

type actionResult struct {
	result status
	logs   string
}

const (
	// NotRun signifies an action not executed
	NotRun status = 0

	// InProgress signifies an action that has been started
	InProgress status = 1

	// Failed signifies an action that failed
	Failed status = 2

	// Succeeded signifies an action that was successful
	Succeeded status = 3
)

type status int
