package deployment

import (
	"fmt"
	"time"

	"github.com/redsailtechnologies/boatswain/pkg/ddd"
)

var runEntityName = "Run"

var startMessageTemplate = `Name:       %s
Type:       %s
Hold:       %s
Start:      %s
------------
`

var completeMessageTemplate = `------------
Status:     %s
Stop:       %s
`

// RunCreated is the event for when a new run is started
type RunCreated struct {
	Timestamp      int64
	UUID           string
	DeploymentUUID string
	Template       *Template
	Trigger        *Trigger
}

// EventType marks this as an event
func (e RunCreated) EventType() string {
	return runEntityName + "Created"
}

// RunStarted signifies when this run is started
type RunStarted struct {
	Timestamp int64
}

// EventType marks this as an event
func (e RunStarted) EventType() string {
	return runEntityName + "Started"
}

// StepStarted is the event for when a particular step is started
type StepStarted struct {
	Timestamp int64
	Name      string
}

// EventType marks this as an event
func (e StepStarted) EventType() string {
	return runEntityName + "StepStarted"
}

// StepCompleted is the event for when a particular step is completed
type StepCompleted struct {
	Timestamp int64
	Status    Status
	Logs      []log
}

// EventType marks this as an event
func (e StepCompleted) EventType() string {
	return runEntityName + "StepCompleted"
}

// RunCompleted signifies when this run is completed
type RunCompleted struct {
	Timestamp int64
	Status    Status
}

// EventType marks this as an event
func (e RunCompleted) EventType() string {
	return runEntityName + "Completed"
}

// Run represents a single execution of a deployment
type Run struct {
	events  []ddd.Event
	version int

	uuid           string
	deploymentUUID string
	trigger        *Trigger
	status         Status

	*Template
}

// FIXME Adamp - we need consistency among these types
type log struct {
	level   logLevel
	message string
}

type logLevel int

const (
	// Debug log level
	Debug logLevel = 0

	// Info log level
	Info logLevel = 1

	// Warn log level
	Warn logLevel = 2

	// Error log level
	Error logLevel = 3
)

// Status represents the outcome of a step
type Status string

const (
	// NotStarted signifies a step not yet executed
	NotStarted Status = "NotStarted"

	// InProgress signifies a step that has been started
	InProgress Status = "In Progress"

	// Failed signifies a step that failed
	Failed Status = "Failed"

	// Succeeded signifies a step that was successful
	Succeeded Status = "Succeeded"

	// Skipped signifies a step that was not run
	Skipped Status = "Skipped"
)

// ReplayRun recreates the run from a series of events
func ReplayRun(events []ddd.Event) *Run {
	r := &Run{}
	for _, event := range events {
		r.on(event)
	}
	return r
}

// CreateRun handles create commands for runs
func CreateRun(uuid, depUUID string, timestamp int64, template *Template, trigger *Trigger) (*Run, error) {
	if uuid == "" {
		return nil, ddd.IDError{}
	}

	r := &Run{}
	r.on(&RunCreated{
		Timestamp:      timestamp,
		UUID:           uuid,
		DeploymentUUID: depUUID,
		Template:       template,
		Trigger:        trigger,
	})
	return r, nil
}

// Start handles start commands for runs
func (r *Run) Start(timestamp int64) error {
	if r.status != NotStarted {
		return ExecutionError{
			Message: "run has already been started",
		}
	}
	r.on(&RunStarted{Timestamp: timestamp})
	return nil
}

// StartStep starts the next step
func (r *Run) StartStep(stepName string, timestamp int64) error {
	if r.status == NotStarted {
		return ExecutionError{
			Message: "run not started",
		}
	}
	for _, step := range *r.Strategy {
		if step.Status == InProgress {
			return ExecutionError{
				Message: "step is still in progress",
				Step:    step.Name,
			}
		} else if step.Name == stepName && step.Status == NotStarted {
			r.on(&StepStarted{
				Timestamp: timestamp,
				Name:      stepName,
			})
			return nil
		}
	}
	return ExecutionError{
		Message: "step not found",
		Step:    stepName,
	}
}

// CompleteStep completes the currently running step
func (r *Run) CompleteStep(status Status, logs []log, timestamp int64) error {
	if r.status == NotStarted && status != Skipped {
		return ExecutionError{
			Message: "run not started",
		}
	}
	for _, step := range *r.Strategy {
		if step.Status == InProgress || (step.Status == NotStarted && status == Skipped) {
			r.on(&StepCompleted{
				Timestamp: timestamp,
				Status:    status,
				Logs:      logs,
			})
			return nil
		}
	}
	return ExecutionError{
		Message: "no step is currently in progress",
	}
}

// Complete handles completion commands for runs
func (r *Run) Complete(status Status, timestamp int64) error {
	r.on(&RunCompleted{
		Timestamp: timestamp,
		Status:    status,
	})
	return nil
}

// UUID gets the run's uuid
func (r *Run) UUID() string {
	return r.uuid
}

// DeploymentUUID gets the run's deployment uuid
func (r *Run) DeploymentUUID() string {
	return r.deploymentUUID
}

// NextStep gets the next step yet to be executed
func (r *Run) NextStep() (*Step, error) {
	for _, step := range *r.Strategy {
		if step.Status == InProgress {
			return nil, ExecutionError{
				Message: "step is still in progress",
				Step:    step.Name,
			}
		} else if step.Status == NotStarted {
			return &step, nil
		}
	}
	return nil, nil
}

// RunVersion gets the version specified in the template for the run (NOTE: not the entity version)
func (r *Run) RunVersion() string {
	return r.Template.Version
}

// Status gets the runs current status
func (r *Run) Status() Status {
	return r.status
}

// Steps gets the steps for this run
func (r *Run) Steps() []Step {
	return *r.Strategy
}

// Events gets the run's event history
func (r *Run) Events() []ddd.Event {
	return r.events
}

// Version gets the run's entity version
func (r *Run) Version() int {
	return r.version
}

func (r *Run) on(event ddd.Event) {
	r.events = append(r.events, event)
	r.version++
	switch e := event.(type) {
	case *RunCreated:
		r.uuid = e.UUID
		r.deploymentUUID = e.DeploymentUUID
		r.trigger = e.Trigger
		r.status = NotStarted
		r.Template = e.Template
		for i := range *r.Strategy {
			(*r.Strategy)[i].Status = NotStarted
		}
	case *RunStarted:
		r.status = InProgress
	case *StepStarted:
		for i, s := range *r.Strategy {
			if s.Name == e.Name && s.Status == NotStarted {
				(*r.Strategy)[i].Status = InProgress
				start := time.Unix(e.Timestamp, 0)
				(*r.Strategy)[i].Logs = append(s.Logs, log{
					level: Info,
					message: fmt.Sprintf(startMessageTemplate,
						s.Name, s.getType(), s.Hold, start.String()),
				})
			}
		}
	case *StepCompleted:
		for i, step := range *r.Strategy {
			if step.Status == InProgress || (step.Status == NotStarted && e.Status == Skipped) {
				(*r.Strategy)[i].Status = e.Status
				stop := time.Unix(e.Timestamp, 0)
				level := Info
				if e.Status == Failed {
					level = Error
				}
				(*r.Strategy)[i].Logs = append((*r.Strategy)[i].Logs, e.Logs...)
				(*r.Strategy)[i].Logs = append((*r.Strategy)[i].Logs, log{
					level: level,
					message: fmt.Sprintf(completeMessageTemplate,
						e.Status, stop.String()),
				})
			}
		}
	case *RunCompleted:
		r.status = e.Status
	}
}

// A ExecutionError is thrown when something goes wrong with execution of a run
type ExecutionError struct {
	Message string
	Step    string
}

func (e ExecutionError) Error() string {
	if e.Step == "" {
		return e.Message
	}
	return e.Step + " " + e.Message
}
