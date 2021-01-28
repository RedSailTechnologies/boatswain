package run

import (
	"bytes"
	"errors"
	"fmt"
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
	"github.com/redsailtechnologies/boatswain/pkg/kube"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/pkg/repo"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
)

var retryCount int = 3

// The Engine is the worker that performs all run steps and tracking
type Engine struct {
	run      *Run
	write    *writeRepository
	clusters *cluster.ReadRepository
	repos    *repo.ReadRepository

	git  git.Agent
	helm helm.Agent
	kube kube.Agent
}

// NewEngine initializes the engine with required dependencies
func NewEngine(r *Run, s storage.Storage, g git.Agent, h helm.Agent, k kube.Agent) (*Engine, error) {
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
		git:      g,
		helm:     h,
		kube:     k,
	}
	return engine, nil
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
	logger.Debug("template for run", "uuid", e.run.UUID(), "template", e.run.template)

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

		var raw []byte
		if step.App.Helm.Command == "install" || step.App.Helm.Command == "upgrade" {
			raw, err = e.downloadChart(app.Helm.Chart, app.Helm.Version, app.Helm.Repo)
			if err != nil {
				e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
				return Failed
			}
		}

		var chart *chart.Chart
		if raw != nil {
			chart, err = loader.LoadArchive(bytes.NewBuffer(raw))
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
			Cluster:   cluster.Name(),
			Endpoint:  cluster.Endpoint(),
			Token:     cluster.Token(),
			Cert:      cluster.Cert(),
			Chart:     chart,
			Values:    vals,
			Wait:      step.App.Helm.Wait,
			Logger:    helmLogger,
		}

		logger.Debug("helm args", "name", args.Name, "namespace", args.Namespace, "values", vals)

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

		logger.Debug("output", "output", output)
		logger.Debug("logs", "logs", logs.String())

		if err != nil {
			logger.Debug("error", "error", err.Error())
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
	return e.helm.GetChart(c, v, endpoint, token)
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
