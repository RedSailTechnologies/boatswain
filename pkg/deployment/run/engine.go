package run

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"time"

	// TODO AdamP - we could consider factoring this into helm.DefaultAgent
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"

	"github.com/redsailtechnologies/boatswain/pkg/cluster"
	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/deployment/template"
	"github.com/redsailtechnologies/boatswain/pkg/git"
	"github.com/redsailtechnologies/boatswain/pkg/helm"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/pkg/repo"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
	"github.com/redsailtechnologies/boatswain/rpc/agent"
)

var retryCount int = 3

// The Engine is the worker that performs all run steps and tracking
type Engine struct {
	run      *Run
	write    *writeRepository
	clusters *cluster.ReadRepository
	repos    *repo.ReadRepository

	agent agent.AgentAction
	git   git.Agent
	repo  repo.Agent
}

// NewEngine initializes the engine with required dependencies
func NewEngine(r *Run, s storage.Storage, a agent.AgentAction, g git.Agent, ra repo.Agent) (*Engine, error) {
	w := newWriteRepository(s)
	if err := w.save(r); err != nil {
		logger.Error("could not save created run", "error", err)
		return nil, err
	}
	engine := &Engine{
		run:      r,
		write:    w,
		clusters: cluster.NewReadRepository(s),
		repos:    repo.NewReadRepository(s),
		agent:    a,
		git:      g,
		repo:     ra,
	}
	return engine, nil
}

// Run starts the engine, and is designed to be run in the background
func (e *Engine) Run() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("recovering from panic", "error", err)
			e.run.AppendLog(err.(error).Error(), Error, ddd.NewTimestamp())
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
		var err error
		app, err := getApp(step.App.Name, e.run.Apps())
		if err != nil {
			e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
			return Failed
		} else if app.Helm == nil {
			e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
			return Failed
		}

		clusters, err := e.clusters.All()
		if err != nil {
			e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
			return Failed
		}
		cluster := getCluster(step.App.Cluster, clusters)
		if cluster == nil {
			e.run.AppendLog("cluster not found", Error, ddd.NewTimestamp())
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

		var vals map[string]interface{}
		if step.App.Helm.Values != nil {
			if step.App.Helm.Values.Library != nil {
				lib := *step.App.Helm.Values.Library
				libRaw, err := e.downloadChart(lib.Chart, lib.Version, lib.Repo)
				if err != nil {
					e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
					return Failed
				}

				libChart, err := loader.LoadArchive(bytes.NewBuffer(libRaw))
				if err != nil {
					e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
					return Failed
				}

				for _, file := range libChart.Files {
					if file.Name == lib.File {
						v, err := chartutil.ReadValues(file.Data)
						if err != nil {
							e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
							return Failed
						}
						vals = v.AsMap()
					}
				}
			}

			if step.App.Helm.Values.Raw != nil {
				// order is important - favor raw values
				vals = mergeVals(*step.App.Helm.Values.Raw, vals)
			}
		}

		args := helm.Args{
			Name:      step.App.Name,
			Namespace: step.App.Namespace,
			Chart:     chart,
			Values:    vals,
			Wait:      step.App.Helm.Wait,
		}
		jsonArgs, err := json.Marshal(args)
		if err != nil {
			e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
			return Failed
		}

		action := &agent.Action{
			Uuid:           ddd.NewUUID(),
			ClusterUuid:    cluster.UUID(),
			ClusterToken:   cluster.Token(),
			ActionType:     agent.ActionType_HELM_ACTION,
			TimeoutSeconds: 600, // FIXME - configurable?
			Args:           jsonArgs,
		}

		var result *agent.Result
		switch step.App.Helm.Command {
		case "install":
			action.Action = string(helm.Install)
		case "rollback":
			action.Action = string(helm.Rollback)
		case "uninstall":
			action.Action = string(helm.Uninstall)
		case "upgrade":
			action.Action = string(helm.Upgrade)
		default:
			e.run.AppendLog("helm command not found", Error, ddd.NewTimestamp())
			return Failed
		}

		result, err = e.agent.Run(context.Background(), action)
		if err != nil {
			e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
			return Failed
		} else if result.Error != "" {
			e.run.AppendLog(result.Error, Error, ddd.NewTimestamp())
			return Failed
		}

		output := ""
		logs := ""
		if action.Action != string(helm.Rollback) && action.Action != string(helm.Uninstall) {
			var r *release.Release
			r, logs, err = helm.ConvertRelease(result.Data)
			if err != nil {
				e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
				return Failed
			}
			output = r.Info.Notes
		} else if action.Action == string(helm.Uninstall) {
			var u *release.UninstallReleaseResponse
			u, logs, err = helm.ConvertUninstallReleaseResponse(result.Data)
			if err != nil {
				e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
				return Failed
			}
			output = u.Info
		} else {
			logs, err = helm.ConvertNone(result.Data)
			if err != nil {
				e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
				return Failed
			}
		}

		e.run.AppendLog(logs, Info, ddd.NewTimestamp())
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
		if err := e.write.save(e.run); err != nil {
			logger.Error("error saving run", "error", err)
		} else {
			return
		}
	}
}

func (e *Engine) downloadChart(c, v, r string) ([]byte, error) {
	if c == "" || v == "" || r == "" {
		return nil, errors.New("cannot download chart without chart, version, and name")
	}
	allRepos, err := e.repos.All()
	if err != nil {
		return nil, err
	}

	endpoint, token, err := getRepoInfo(r, allRepos)
	if err != nil {
		return nil, err
	}
	return e.repo.GetChart(c, v, endpoint, token)
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

func getCluster(c string, l []*cluster.Cluster) *cluster.Cluster {
	for _, cluster := range l {
		if cluster.Name() == c {
			return cluster
		}
	}
	return nil
}

func getRepoInfo(r string, l []*repo.Repo) (string, string, error) {
	for _, repo := range l {
		if repo.Name() == r {
			return repo.Endpoint(), repo.Token(), nil
		}
	}
	return "", "", errors.New("repo not found")
}

func mergeVals(one, two map[string]interface{}) map[string]interface{} {
	if one == nil {
		one = make(map[string]interface{})
	}
	c := &chart.Chart{}
	c.Values = one

	v, err := chartutil.CoalesceValues(c, two)
	if err != nil {
		return nil
	}
	return v.AsMap()
}
