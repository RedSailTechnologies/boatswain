package deployment

import (
	"bytes"
	"errors"

	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/git"
	"github.com/redsailtechnologies/boatswain/pkg/helm"
	"github.com/redsailtechnologies/boatswain/pkg/kube"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/rpc/cluster"
	"github.com/redsailtechnologies/boatswain/rpc/repo"
)

var retryCount int = 3

type runAgents struct {
	git  git.Agent
	helm helm.Agent
	kube kube.Agent
}

type runEntities struct {
	clusters []*cluster.ClusterRead
	repos    []*repo.RepoRead
}

// The RunEngine is the worker that performs all run steps and tracking
type RunEngine struct {
	run      *Run
	rr       *RunRepository
	agents   runAgents
	entities runEntities
}

// NewRunEngine initializes the engine with required dependencies
func NewRunEngine(r *Run, rr *RunRepository, a runAgents, e runEntities) *RunEngine {
	// TODO AdamP - we could filter out only what we need for repos and clusters
	engine := &RunEngine{
		run:      r,
		rr:       rr,
		agents:   a,
		entities: e,
	}
	return engine
}

// Run runs the engine, and is designed to be run in the background
func (e *RunEngine) Run() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("recovering from panic", "error", err)
			e.finalize(Failed)
		}
	}()

	logger.Info("starting run", "uuid", e.run.UUID())
	err := e.run.Start(ddd.NewTimestamp())
	if err != nil {
		logger.Warn("error starting run", "error", err)
		e.finalize(Failed)
		return
	}
	e.persist()

	var lastStatus Status = Succeeded
	var overallStatus = Succeeded
	for true {
		step, err := e.run.NextStep()
		if err != nil {
			logger.Warn("next step called with step in progress", "error", err)
			e.finalize(Failed)
			return
		}

		// step is nil if there are no more steps to be executed
		if step == nil && err == nil {
			e.finalize(overallStatus)
			return
		}

		if step.shouldExecute(lastStatus) {
			e.run.StartStep(step.Name, ddd.NewTimestamp())
			e.persist()

			status, logs := e.executeStep(step)

			e.run.CompleteStep(status, logs, ddd.NewTimestamp())
			err = e.rr.Save(e.run)
			if err != nil {
				logger.Error("could not persist run", "error", err)
				e.finalize(Failed)
				return
			}

			if status == Failed {
				overallStatus = Failed
			}

			lastStatus = status
		} else {
			e.run.CompleteStep(Skipped, "", ddd.NewTimestamp())
			e.persist()
		}
	}
}

func (e *RunEngine) executeStep(step *Step) (Status, Log) {
	if step.App != nil {
		s, l := e.executeActionStep(step)
		return s, Log(l)
	} else if step.Test != nil {
		return Skipped, "not implemented"
	} else if step.Approval != nil {
		return Skipped, "not implemented"
	} else if step.Trigger != nil {
		return Skipped, "not implemented"
	}
	return Failed, ""
}

func (e *RunEngine) executeActionStep(step *Step) (Status, string) {
	if step.App.Helm != nil {
		app, err := getApp(step.App.Name, *e.run.Apps)
		if err != nil {
			return Failed, err.Error()
		} else if app.Helm == nil {
			return Failed, "defined app is not a helm app"
		}

		endpoint, token, err := getClusterInfo(step.App.Cluster, e.entities.clusters)
		if err != nil {
			return Failed, err.Error()
		}

		var chart []byte
		if step.App.Helm.Command == "install" || step.App.Helm.Command == "upgrade" {
			chart, err = e.downloadChart(app.Helm.Chart, app.Helm.Version, app.Helm.Repo)
			if err != nil {
				return Failed, err.Error()
			}
		}

		// FIXME get raw and library values and merge them
		vals := make(map[string]interface{})

		log := &bytes.Buffer{}

		args := helm.Args{
			Name:      step.App.Name,
			Namespace: step.App.Namespace,
			Endpoint:  endpoint,
			Token:     token,
			Chart:     chart,
			Values:    vals, // FIXME
			Wait:      step.App.Helm.Wait,
			Logger:    log,
		}

		switch step.App.Helm.Command {
		case "install":
			_, err = e.agents.helm.Install(args)
		case "rollback":
			err = e.agents.helm.Rollback(step.App.Helm.Version, args)
		case "uninstall":
			_, err = e.agents.helm.Uninstall(args)
		case "upgrade":
			_, err = e.agents.helm.Upgrade(args)
		default:
			return Failed, "helm command not found"
		}

		if err != nil {
			return Failed, err.Error()
		}
		return Succeeded, log.String()
	}
	return Failed, "only helm apps are currently supported"
}

func (e *RunEngine) finalize(status Status) {
	logger.Info("completing run", "uuid", e.run.UUID(), "status", status)
	for i := 0; i < retryCount; i++ {
		e.run.Complete(status, ddd.NewTimestamp())
		if err := e.rr.Save(e.run); err == nil {
			return
		}
	}
}

func (e *RunEngine) downloadChart(c, v, r string) ([]byte, error) {
	if c == "" || v == "" || r == "" {
		return nil, errors.New("cannot download chart without chart, version, and name")
	}
	repo, err := getRepoInfo(r, e.entities.repos)
	if err != nil {
		return nil, err
	}
	return e.agents.helm.GetChart(c, v, repo)
}

func (e *RunEngine) persist() {
	err := e.rr.Save(e.run)
	if err != nil {
		logger.Error("could not persist run", "error", err)
		e.finalize(Failed)
	}
}

func (s *Step) shouldExecute(last Status) bool {
	return s.Condition == "always" || s.Condition == "any" ||
		(s.Condition == "succeeded" && last == Succeeded) ||
		(s.Condition == "failed" && last == Failed)
}

func getApp(name string, apps []App) (*App, error) {
	for _, app := range apps {
		if app.Name == name {
			return &app, nil
		}
	}
	return nil, errors.New("app definition not found")
}

func getClusterInfo(c string, l []*cluster.ClusterRead) (string, string, error) {
	for _, cluster := range l {
		if cluster.Name == c {
			return cluster.Endpoint, cluster.Token, nil
		}
	}
	return "", "", errors.New("cluster not found")
}

func getRepoInfo(r string, l []*repo.RepoRead) (string, error) {
	for _, repo := range l {
		if repo.Name == r {
			return repo.Endpoint, nil
		}
	}
	return "", errors.New("repo not found")
}
