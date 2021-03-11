package run

import (
	"time"

	"github.com/redsailtechnologies/boatswain/pkg/cluster"
	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/deployment/approval"
	"github.com/redsailtechnologies/boatswain/pkg/deployment/template"
	"github.com/redsailtechnologies/boatswain/pkg/git"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/pkg/repo"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
	"github.com/redsailtechnologies/boatswain/rpc/agent"
)

var retryCount int = 3

// The Engine is the worker that performs all run steps and tracking
type Engine struct {
	run      *Run
	statuses *statuses

	write     *writeRepository
	clusters  *cluster.ReadRepository
	repos     *repo.ReadRepository
	aprvRead  *approval.ReadRepository
	aprvWrite *approval.WriteRepository

	agent agent.AgentAction
	git   git.Agent
	repo  repo.Agent

	trigger func(string, string, []byte) (string, error)
}

// NewEngine initializes the engine with required dependencies
func NewEngine(r *Run, s storage.Storage, a agent.AgentAction, g git.Agent, ra repo.Agent, t func(string, string, []byte) (string, error)) (*Engine, error) {
	w := newWriteRepository(s)
	if err := w.save(r); err != nil {
		logger.Error("could not save created run", "error", err)
		return nil, err
	}
	engine := &Engine{
		run:       r,
		statuses:  &statuses{},
		write:     w,
		clusters:  cluster.NewReadRepository(s),
		repos:     repo.NewReadRepository(s),
		aprvRead:  approval.NewReadRepository(s),
		aprvWrite: approval.NewWriteRepository(s),
		agent:     a,
		git:       g,
		repo:      ra,
		trigger:   t,
	}
	return engine, nil
}

// Start the engine-designed to be run in the background
func (e *Engine) Start() {
	defer e.recover()

	logger.Info("starting run", "uuid", e.run.UUID(), "start", ddd.NewTimestamp())
	err := e.run.Start(ddd.NewTimestamp())
	if err != nil {
		logger.Warn("error starting run", "error", err)
		e.finalize(Failed)
		return
	}
	e.persist()
	e.startLoop(Succeeded, Succeeded)
}

// Resume continues a run after a pause
func (e *Engine) Resume() {
	defer e.recover()

	steps := e.run.Steps()
	last := steps[len(steps)-1].Status
	overall := Succeeded
	for _, step := range steps {
		if step.Status == Failed {
			overall = Failed
		}
	}
	logger.Info("continuing run", "uuid", e.run.UUID(), "start", ddd.NewTimestamp())
	e.startLoop(last, overall)
}

func (e *Engine) startLoop(lastStatus, overallStatus Status) {
	defer e.recover()

	status := lastStatus
	overall := overallStatus
	looping := true

	for looping {
		step := e.run.CurrentTemplate()

		// step is nil if there are no more steps to be executed
		if step == nil {
			e.finalize(overall)
			return
		}

		if shouldExecute(step, status) {
			if !e.run.Paused() {
				e.run.StartStep(step.Name, ddd.NewTimestamp())
				e.persist()

				// the hold is executed at the beginning of the step before any logic is executed
				// so here its safe to run this only if we're starting the step, not if its been paused
				if step.Hold != "" {
					e.executeHold(step.Hold)
				}
			}

			status = e.executeStep(step)
			e.run.SetStatus(status)
			if status == Failed || status == Succeeded {
				e.run.CompleteStep(ddd.NewTimestamp())
			}
			e.persist()

			if status == Failed {
				overall = Failed
			}

			if e.run.Paused() {
				looping = false
			}
		} else {
			e.run.SkipStep(step.Name, "conditions not met", ddd.NewTimestamp())
			e.persist()
		}
	}
}

func (e *Engine) executeStep(step *template.Step) Status {
	if step.Helm != nil {
		return e.executeHelmStep(step)
	} else if step.Approval != nil {
		return e.executeApprovalStep(step)
	} else if step.Trigger != nil {
		return e.executeTriggerStep(step)
	}
	e.run.AppendLog("step has nothing to execute", Error, ddd.NewTimestamp())
	return Failed
}

func (e *Engine) finalize(s Status) {
	logger.Info("completing run", "uuid", e.run.UUID(), "status", s)
	for i := 0; i < retryCount; i++ {
		if err := e.run.Complete(s, ddd.NewTimestamp()); err != nil {
			logger.Error("error completing run", "error", err)
			continue
		}
		if err := e.write.save(e.run); err != nil {
			logger.Error("error saving run", "error", err)
		} else {
			return
		}
	}
}

func (e *Engine) persist() {
	err := e.write.save(e.run)
	if err != nil {
		logger.Error("could not persist run", "error", err)
		e.finalize(Failed)
	}
}

func (e *Engine) executeHold(h string) error {
	d, err := time.ParseDuration(h)
	if err != nil {
		return err
	}
	time.Sleep(d)
	return nil
}

func (e *Engine) recover() {
	if err := recover(); err != nil {
		logger.Error("recovering from panic", "error", err)
		e.run.AppendLog(err.(error).Error(), Error, ddd.NewTimestamp())
		e.finalize(Failed)
	}
}

func shouldExecute(s *template.Step, last Status) bool {
	return s.Condition == "always" || s.Condition == "any" ||
		(s.Condition == "succeeded" && last == Succeeded) ||
		(s.Condition == "failed" && last == Failed)
}
