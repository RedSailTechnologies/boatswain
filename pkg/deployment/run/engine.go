package run

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"

	"github.com/redsailtechnologies/boatswain/pkg/cluster"
	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/deployment/approval"
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
	run       *Run
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
	if step.App != nil || step.Test != nil {
		return e.executeActionStep(step)
	} else if step.Approval != nil {
		return e.executeApprovalStep(step)
	} else if step.Trigger != nil {
		return e.executeTriggerStep(step)
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
		case "test":
			action.Action = string(helm.Test)
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

func (e *Engine) executeTriggerStep(step *template.Step) Status {
	b, err := yaml.Marshal(step.Trigger.Arguments)
	if err != nil {
		e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
		return Failed
	}

	res, err := e.trigger(step.Trigger.Name, step.Trigger.Deployment, b)
	if err != nil {
		e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
		return Failed
	}
	e.run.AppendLog(fmt.Sprintf("triggered run with uuid %s", res), Info, ddd.NewTimestamp())
	return Succeeded
}

func (e *Engine) executeApprovalStep(step *template.Step) Status {
	approvals, err := e.aprvRead.All()
	if err != nil {
		e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
		return Failed
	}

	var appr *approval.Approval
	for _, a := range approvals {
		if a.RunUUID() == e.run.UUID() {
			appr = a
		}
	}
	if appr == nil {
		appr, err = approval.Create(ddd.NewUUID(), e.run.UUID(), step.Approval.Users, step.Approval.Roles)
		if err != nil {
			e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
			return Failed
		}

		err = e.aprvWrite.Save(appr)
		if err != nil {
			e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
			return Failed
		}
		e.run.AppendLog("awaiting approval", Info, ddd.NewTimestamp())
		return AwaitingApproval
	} else if appr.Completed() {
		appr.Destroy()
		err := e.aprvWrite.Save(appr)
		if err != nil {
			logger.Warn("step could not be approved", "error", err, "run", e.run.UUID())
			return AwaitingApproval
		}
		approver, err := appr.Approver()
		if err != nil {
			e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
			return Failed
		}

		if appr.Overridden() {
			e.run.AppendLog(fmt.Sprintf("step overriden by %s", approver.Name), Info, ddd.NewTimestamp())
			return Succeeded
		} else if appr.Approved() {
			e.run.AppendLog(fmt.Sprintf("step approved by %s", approver.Name), Info, ddd.NewTimestamp())
			return Succeeded
		}
		e.run.AppendLog(fmt.Sprintf("step rejected by %s", approver.Name), Info, ddd.NewTimestamp())
		return Failed
	}
	e.run.AppendLog("approval step not found or could not be created", Error, ddd.NewTimestamp())
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
