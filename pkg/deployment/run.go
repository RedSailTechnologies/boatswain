package deployment

import "github.com/redsailtechnologies/boatswain/pkg/ddd"

var runEntityName = "Run"

// RunCreated is the event for when a new run is started
type RunCreated struct {
	Timestamp int64
	UUID      string
	Template  *Template
	Trigger   *Trigger
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
	Logs      Logs
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

	uuid    string
	trigger *Trigger
	status  Status

	*Template
}

// Logs represents the logs for a step
type Logs []string

// Status represents the outcome of a step
type Status int

const (
	// NotStarted signifies a step not yet executed
	NotStarted Status = 0

	// InProgress signifies a step that has been started
	InProgress Status = 1

	// Failed signifies a step that failed
	Failed Status = 2

	// Succeeded signifies a step that was successful
	Succeeded Status = 3

	// Skipped signifies a step that was not run
	Skipped Status = 4
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
func CreateRun(uuid string, timestamp int64, template *Template, trigger *Trigger) (*Run, error) {
	if uuid == "" {
		return nil, ddd.IDError{}
	}

	r := &Run{}
	r.on(&RunCreated{
		Timestamp: timestamp,
		UUID:      uuid,
		Template:  template,
		Trigger:   trigger,
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
func (r *Run) CompleteStep(status Status, logs Logs, timestamp int64) error {
	if r.status == NotStarted && status != Skipped {
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
		} else if step.Status == InProgress {
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

// TODO - getters

// UUID gets the run's uuid
func (r *Run) UUID() string {
	return r.uuid
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

// Status gets the runs current status
func (r *Run) Status() Status {
	return r.status
}

// Events gets the run's event history
func (r *Run) Events() []ddd.Event {
	return r.events
}

func (r *Run) on(event ddd.Event) {
	r.events = append(r.events, event)
	r.version++
	switch e := event.(type) {
	case *RunCreated:
		r.uuid = e.UUID
		r.trigger = e.Trigger
		r.status = NotStarted
		r.Template = e.Template
	case *StepStarted:
		if r.status == NotStarted {
			r.status = InProgress
		}
		for _, step := range *r.Strategy {
			if step.Name == e.Name && step.Status == NotStarted {
				step.Status = InProgress
			}
		}
	case *StepCompleted:
		for _, step := range *r.Strategy {
			if step.Status == InProgress {
				step.Status = e.Status
				step.Logs = e.Logs
				return
			}
		}
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
