package run

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/deployment/template"
	"github.com/redsailtechnologies/boatswain/pkg/git"
	"github.com/redsailtechnologies/boatswain/pkg/helm"
	"github.com/redsailtechnologies/boatswain/pkg/kube"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/rpc/cluster"
	"github.com/redsailtechnologies/boatswain/rpc/repo"
	"helm.sh/helm/v3/pkg/release"
)

var retryCount int = 3

// The Engine is the worker that performs all run steps and tracking
type Engine struct {
	run  *Run
	repo *Repository

	git  git.Agent
	helm helm.Agent
	kube kube.Agent

	clusters []*cluster.ClusterRead
	repos    []*repo.RepoRead
}

// NewEngine initializes the engine with required dependencies
func NewEngine(r *Run, repo *Repository, g git.Agent, h helm.Agent, k kube.Agent, cl []*cluster.ClusterRead, re []*repo.RepoRead) *Engine {
	// TODO AdamP - we could filter out only what we need for repos and clusters
	engine := &Engine{
		run:      r,
		repo:     repo,
		git:      g,
		helm:     h,
		kube:     k,
		clusters: cl,
		repos:    re,
	}
	return engine
}

// Run starts the engine, and is designed to be run in the background
func (e *Engine) Run() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("recovering from panic", "error", err)
			e.finalize(Failed)
		}
	}()

	logger.Info("starting run", "uuid", e.run.UUID(), "start", ddd.NewTimestamp())

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
		step := e.run.CurrentTemplate()

		// step is nil if there are no more steps to be executed
		if step == nil {
			e.finalize(overallStatus)
			return
		}

		if shouldExecute(step, lastStatus) {
			e.run.StartStep(step.Name, ddd.NewTimestamp())
			e.persist()

			if step.Hold != "" {
				e.executeHold(step.Hold)
			}

			status := e.executeStep(step)

			e.run.CompleteStep(status, ddd.NewTimestamp())
			e.persist()

			if status == Failed {
				overallStatus = Failed
			}

			lastStatus = status
		} else {
			e.run.SkipStep(step.Name, "conditions not met", ddd.NewTimestamp())
			e.persist()
		}
	}
}

func (e *Engine) executeStep(step *template.Step) Status {
	if step.App != nil {
		return e.executeActionStep(step)
	} else if step.Test != nil {
		e.run.AppendLog("step type not implemented", Error, ddd.NewTimestamp())
		return Skipped
	} else if step.Approval != nil {
		e.run.AppendLog("step type not implemented", Error, ddd.NewTimestamp())
		return Skipped
	} else if step.Trigger != nil {
		e.run.AppendLog("step type not implemented", Error, ddd.NewTimestamp())
		return Skipped
	}
	e.run.AppendLog("step has nothing to execute", Error, ddd.NewTimestamp())
	return Failed
}

func (e *Engine) executeActionStep(step *template.Step) Status {
	if step.App.Helm != nil {
		app, err := getApp(step.App.Name, e.run.Apps())
		if err != nil {
			e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
			return Failed
		} else if app.Helm == nil {
			e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
			return Failed
		}

		endpoint, token, err := getClusterInfo(step.App.Cluster, e.clusters)
		if err != nil {
			e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
			return Failed
		}

		var chart []byte
		if step.App.Helm.Command == "install" || step.App.Helm.Command == "upgrade" {
			chart, err = e.downloadChart(app.Helm.Chart, app.Helm.Version, app.Helm.Repo)
			if err != nil {
				e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
				return Failed
			}
		}

		// FIXME get library values and merge them
		vals := make(map[string]interface{})
		if step.App.Helm.Values != nil && step.App.Helm.Values.Raw != nil {
			vals = *step.App.Helm.Values.Raw
		}

		logs := &bytes.Buffer{}
		helmLogger := func(t string, a ...interface{}) {
			str := fmt.Sprintf(t, a...)
			if str[len(str)-1] != '\n' {
				str = str + "\n"
			}
			logs.Write([]byte(str))
		}

		args := helm.Args{
			Name:      step.App.Name,
			Namespace: step.App.Namespace,
			Endpoint:  endpoint,
			Token:     token,
			Chart:     chart,
			Values:    vals, // FIXME
			Wait:      step.App.Helm.Wait,
			Logger:    helmLogger,
		}

		output := ""
		switch step.App.Helm.Command {
		case "install":
			var r *release.Release
			r, err = e.helm.Install(args)
			if r != nil {
				output = r.Info.Notes
			}
		case "rollback":
			err = e.helm.Rollback(step.App.Helm.Version, args)
		case "uninstall":
			var r *release.UninstallReleaseResponse
			r, err = e.helm.Uninstall(args)
			output = r.Info
		case "upgrade":
			var r *release.Release
			r, err = e.helm.Upgrade(args)
			if r != nil {
				output = r.Info.Notes
			}
		default:
			e.run.AppendLog("helm command not found", Error, ddd.NewTimestamp())
			return Failed
		}

		if err != nil {
			e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
			return Failed
		}

		e.run.AppendLog(logs.String(), Info, ddd.NewTimestamp())
		e.run.AppendLog(output, Info, ddd.NewTimestamp())
		return Succeeded
	}

	e.run.AppendLog("only helm apps are currently supported", Error, ddd.NewTimestamp())
	return Failed
}

func (e *Engine) finalize(status Status) {
	logger.Info("completing run", "uuid", e.run.UUID(), "status", status)
	for i := 0; i < retryCount; i++ {
		if err := e.run.Complete(status, ddd.NewTimestamp()); err != nil {
			logger.Error("error completing run", "error", err)
			continue
		}
		if err := e.repo.Save(e.run); err == nil {
			return
		} else {
			logger.Error("error saving run", "error", err)
		}
	}
}

func (e *Engine) downloadChart(c, v, r string) ([]byte, error) {
	if c == "" || v == "" || r == "" {
		return nil, errors.New("cannot download chart without chart, version, and name")
	}
	repo, err := getRepoInfo(r, e.repos)
	if err != nil {
		return nil, err
	}
	return e.helm.GetChart(c, v, repo)
}

func (e *Engine) persist() {
	err := e.repo.Save(e.run)
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

func shouldExecute(s *template.Step, last Status) bool {
	return s.Condition == "always" || s.Condition == "any" ||
		(s.Condition == "succeeded" && last == Succeeded) ||
		(s.Condition == "failed" && last == Failed)
}

func getApp(name string, apps []template.App) (*template.App, error) {
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
