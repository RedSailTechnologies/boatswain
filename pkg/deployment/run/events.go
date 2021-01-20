package run

import (
	"github.com/redsailtechnologies/boatswain/pkg/deployment/template"
	"github.com/redsailtechnologies/boatswain/pkg/deployment/trigger"
)

// Created is the event for when a new run is started
type Created struct {
	UUID           string
	DeploymentUUID string
	Template       *template.Template
	Trigger        *trigger.Trigger
}

// EventType marks this as an event
func (e Created) EventType() string {
	return entityName + "Created"
}

// Started signifies when this run is started
type Started struct {
	Timestamp int64
}

// EventType marks this as an event
func (e Started) EventType() string {
	return entityName + "Started"
}

// StepStarted is the event for when a particular step is started
type StepStarted struct {
	Timestamp int64
	Name      string
}

// EventType marks this as an event
func (e StepStarted) EventType() string {
	return entityName + "StepStarted"
}

// AppendLog is the event for when we append a logging event
type AppendLog struct {
	Timestamp int64
	Level     LogLevel
	Message   string
}

// EventType marks this as an event
func (e AppendLog) EventType() string {
	return entityName + "AppendLog"
}

// StepCompleted is the event for when a particular step is completed
type StepCompleted struct {
	Timestamp int64
	Status    Status
}

// EventType marks this as an event
func (e StepCompleted) EventType() string {
	return entityName + "StepCompleted"
}

// StepSkipped is the event for when a particular step is skipped
type StepSkipped struct {
	Timestamp int64
	Name      string
	Reason    string
}

// EventType marks this as an event
func (e StepSkipped) EventType() string {
	return entityName + "StepSkipped"
}

// Completed signifies when this run is completed
type Completed struct {
	Timestamp int64
	Status    Status
}

// EventType marks this as an event
func (e Completed) EventType() string {
	return entityName + "Completed"
}
