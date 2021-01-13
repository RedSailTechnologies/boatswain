package deployment

import (
	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/rpc/cluster"
	"github.com/redsailtechnologies/boatswain/rpc/repo"
)

var retryCount int = 3

// The RunEngine is the worker that performs all run steps and tracking
type RunEngine struct {
	run *Run
	rr  *RunRepository

	clusters []*cluster.ClusterRead
	repos    []*repo.RepoRead
}

// NewRunEngine initializes the engine with required dependencies
func NewRunEngine(r *Run, rr *RunRepository, cl []*cluster.ClusterRead, rep []*repo.RepoRead) *RunEngine {
	// TODO AdamP - we could filter out only what we need for repos and clusters
	e := &RunEngine{
		run: r,
		rr:  rr,

		clusters: cl,
		repos:    rep,
	}
	return e
}

// Run runs the engine, and is designed to be run in the background
func (e *RunEngine) Run() {
	// TODO - panic handler so we always catch something

	// start the run
	logger.Info("starting run", "uuid", e.run.UUID())
	err := e.run.Start(ddd.NewTimestamp())
	if err != nil {
		logger.Warn("error starting run", "error", err)
		e.killExecution(Failed)
	}
	e.persist()

	var lastStatus Status = Succeeded
	for true {
		step, err := e.run.NextStep()
		if err != nil {
			logger.Warn("next step called with step in progress", "error", err)
			e.killExecution(Failed)
		}

		if step == nil && err == nil {
			// TODO - kill execution, getting the correct status
		}

		if step.shouldExecute(lastStatus) {
			e.run.StartStep(step.Name, ddd.NewTimestamp())
			e.persist()

			status, logs := e.executeStep(step)

			e.run.CompleteStep(status, logs, ddd.NewTimestamp())
			err = e.rr.Save(e.run)
			if err != nil {
				logger.Error("could not persist run", "error", err)
				e.killExecution(Failed)
			}

			lastStatus = status
		} else {
			e.run.CompleteStep(Skipped, nil, ddd.NewTimestamp())
		}

	}
}

func (e *RunEngine) executeStep(step *Step) (Status, Logs) {
	// switches for:
	// helm install
	// helm upgrade
	// helm uninstall
	// helm rollback
	// helm test
	// approval (do we need a channel?)
	// trigger (another channel?)
	return nil, nil
}

// TODO - add an error param here?
func (e *RunEngine) killExecution(status Status) {
	logger.Info("completing run", "uuid", e.run.UUID(), "status", e.run.Status())
	for i := 0; i < retryCount; i++ {
		e.run.Complete(status, ddd.NewTimestamp())
		if err := e.rr.Save(e.run); err == nil {
			return
		}
	}
}

func (e *RunEngine) persist() {
	err := e.rr.Save(e.run)
	if err != nil {
		logger.Error("could not persist run", "error", err)
		e.killExecution(Failed)
	}
}

func (s *Step) shouldExecute(last Status) bool {
	return s.Condition == "always" || s.Condition == "any" ||
		(s.Condition == "succeeded" && last == Succeeded) ||
		(s.Condition == "failed" && last == Failed)
}
