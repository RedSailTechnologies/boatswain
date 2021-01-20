package run

import (
	"fmt"

	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/deployment/template"
	"github.com/redsailtechnologies/boatswain/pkg/deployment/trigger"
)

var entityName = "Run"

var startMessageTemplate = `Step Started
Name:       %s
Type:       %s
Hold:       %s`

var completeMessageTemplate = `Step Completed
Status:     %s`

// Run represents a single execution of a deployment
type Run struct {
	events  []ddd.Event
	version int

	uuid     string
	deployID string
	trigger  *trigger.Trigger
	template *template.Template

	status  Status
	start   int64
	stop    int64
	current int
	steps   []Step
}

// Replay recreates the run from a series of events
func Replay(events []ddd.Event) *Run {
	r := &Run{}
	for _, event := range events {
		r.on(event)
	}
	return r
}

// Create handles create commands for runs
func Create(uuid, depUUID string, tpl *template.Template, trig *trigger.Trigger) (*Run, error) {
	if uuid == "" {
		return nil, ddd.IDError{}
	}
	if depUUID == "" {
		return nil, RuntimeError{m: "deployment uuid cannot be empty"}
	}
	if tpl == nil {
		return nil, RuntimeError{m: "template cannot be nil"}
	}

	r := &Run{}
	r.on(&Created{
		UUID:           uuid,
		DeploymentUUID: depUUID,
		Template:       tpl,
		Trigger:        trig,
	})
	return r, nil
}

// Start handles start commands for runs
func (r *Run) Start(timestamp int64) error {
	if r.status != NotStarted {
		return RuntimeError{m: "run has already been started"}
	}
	r.on(&Started{Timestamp: timestamp})
	return nil
}

// StartStep starts the next step
func (r *Run) StartStep(name string, timestamp int64) error {
	if name == "" {
		return RuntimeError{m: "step name cannot be empty"}
	}
	r.on(&StepStarted{
		Timestamp: timestamp,
		Name:      name,
	})
	return nil
}

// AppendLog appends the log message to the current step
func (r *Run) AppendLog(message string, level LogLevel, timestamp int64) error {
	if message == "" {
		return RuntimeError{m: "cannot log an empty message"}
	}
	r.on(&AppendLog{
		Timestamp: timestamp,
		Level:     level,
		Message:   message,
	})
	return nil
}

// CompleteStep completes the currently running step
func (r *Run) CompleteStep(status Status, timestamp int64) error {
	if status == NotStarted || status == InProgress {
		return RuntimeError{m: "to complete a step status must be terminal"}
	}
	r.on(&StepCompleted{
		Timestamp: timestamp,
		Status:    status,
	})
	return nil
}

// SkipStep skips the step completely
func (r *Run) SkipStep(name, reason string, timestamp int64) error {
	if name == "" {
		return RuntimeError{m: "step name cannot be empty"}
	}
	if reason == "" {
		return RuntimeError{m: "reason required when skipping a step"}
	}
	r.on(&StepSkipped{
		Timestamp: timestamp,
		Name:      name,
		Reason:    reason,
	})
	return nil
}

// Complete handles completion commands for runs
func (r *Run) Complete(status Status, timestamp int64) error {
	r.on(&Completed{
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
	return r.deployID
}

// CurrentTemplate gets the current step to execute
func (r *Run) CurrentTemplate() *template.Step {
	if r.current >= len(*r.template.Strategy) {
		return nil
	}
	return &(*r.template.Strategy)[r.current]
}

// RunVersion gets the version specified in the template for the run (NOTE: not the entity version)
func (r *Run) RunVersion() string {
	return r.template.Version
}

// StartTime gets the time this run started
func (r *Run) StartTime() int64 {
	return r.start
}

// StopTime gets the time this run was completed
func (r *Run) StopTime() int64 {
	return r.stop
}

// Status gets the runs current status
func (r *Run) Status() Status {
	return r.status
}

// Steps gets the steps for this run
func (r *Run) Steps() []Step {
	cp := make([]Step, len(r.steps))
	copy(cp, r.steps)
	return cp
}

// Apps gets the apps for this run
func (r *Run) Apps() []template.App {
	cp := make([]template.App, len(*r.template.Apps))
	copy(cp, *r.template.Apps)
	return cp
}

// Events gets the run's event history
func (r *Run) Events() []ddd.Event {
	cp := make([]ddd.Event, len(r.events))
	copy(cp, r.events)
	return cp
}

// Version gets the run's entity version
func (r *Run) Version() int {
	return r.version
}

func (r *Run) on(event ddd.Event) {
	r.events = append(r.events, event)
	r.version++
	switch e := event.(type) {
	case *Created:
		r.uuid = e.UUID
		r.deployID = e.DeploymentUUID
		r.trigger = e.Trigger
		r.template = e.Template
		r.status = NotStarted
		r.current = 0 // not 100% necessary but more explicit
	case *Started:
		r.start = e.Timestamp
		r.status = InProgress
	case *StepStarted:
		s := Step{
			Name:   e.Name,
			Status: InProgress,
			Start:  e.Timestamp,
		}
		t := (*r.template.Strategy)[r.current]
		s.log(fmt.Sprintf(startMessageTemplate, t.Name, getTemplateType(t), t.Hold), Info, e.Timestamp)
		r.steps = append(r.steps, s)
	case *AppendLog:
		r.steps[r.current].log(e.Message, e.Level, e.Timestamp)
	case *StepCompleted:
		s := &r.steps[r.current]
		s.Status = e.Status
		s.Stop = e.Timestamp
		s.log(fmt.Sprintf(completeMessageTemplate, e.Status), Info, e.Timestamp)
		r.current++
	case *StepSkipped:
		s := Step{
			Name:   e.Name,
			Status: Skipped,
			Start:  e.Timestamp,
			Stop:   e.Timestamp,
		}
		s.log("step skipped: "+e.Reason, Info, e.Timestamp)
		r.steps = append(r.steps, s)
		r.current++
	case *Completed:
		r.status = e.Status
		r.stop = e.Timestamp
	}
}

func getTemplateType(s template.Step) string {
	if s.App != nil {
		return "App"
	} else if s.Test != nil {
		return "Test"
	} else if s.Approval != nil {
		return "Approval"
	} else if s.Trigger != nil {
		return "Trigger"
	} else {
		return ""
	}
}
