package deployment

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/git"
	"github.com/redsailtechnologies/boatswain/pkg/helm"
	"github.com/redsailtechnologies/boatswain/pkg/kube"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/rpc/cluster"
	"github.com/redsailtechnologies/boatswain/rpc/repo"
	"helm.sh/helm/v3/pkg/release"
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

// Run starts the engine, and is designed to be run in the background
func (e *RunEngine) Run() {
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
		step, err := e.run.NextStep()
		if err != nil {
			logger.Warn("error getting the next step", "error", err)
			e.finalize(Failed)
			return
		}

		// step is nil if there are no more steps to be executed
		if step == nil && err == nil {
			e.finalize(overallStatus)
			return
		}

		if step.shouldExecute(lastStatus) {
			if step.Hold != "" {
				e.executeHold(step.Hold)
			}

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
			skipLog := []log{
				log{
					level:   Info,
					message: "step skipped",
				},
			}
			e.run.CompleteStep(Skipped, skipLog, ddd.NewTimestamp())
			e.persist()
		}
	}
}

func (e *RunEngine) executeStep(step *Step) (Status, []log) {
	logs := make([]log, 0)
	if step.App != nil {
		return e.executeActionStep(step, logs)
	} else if step.Test != nil {
		return Skipped, []log{
			log{
				level:   Info,
				message: "not implemented",
			},
		}
	} else if step.Approval != nil {
		return Skipped, []log{
			log{
				level:   Info,
				message: "not implemented",
			},
		}
	} else if step.Trigger != nil {
		return Skipped, []log{
			log{
				level:   Info,
				message: "not implemented",
			},
		}
	}
	return Failed, []log{
		log{
			level:   Info,
			message: "step does not have anything to execute",
		},
	}
}

func (e *RunEngine) executeActionStep(step *Step, l []log) (Status, []log) {
	if step.App.Helm != nil {
		app, err := getApp(step.App.Name, *e.run.Apps)
		if err != nil {
			return Failed, append(l, log{
				level:   Error,
				message: err.Error(),
			})
		} else if app.Helm == nil {
			return Failed, append(l, log{
				level:   Error,
				message: "defined app is not a helm app",
			})
		}

		endpoint, token, err := getClusterInfo(step.App.Cluster, e.entities.clusters)
		if err != nil {
			return Failed, append(l, log{
				level:   Error,
				message: err.Error(),
			})
		}

		var chart []byte
		if step.App.Helm.Command == "install" || step.App.Helm.Command == "upgrade" {
			chart, err = e.downloadChart(app.Helm.Chart, app.Helm.Version, app.Helm.Repo)
			if err != nil {
				return Failed, append(l, log{
					level:   Error,
					message: err.Error(),
				})
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
			r, err = e.agents.helm.Install(args)
			if r != nil {
				output = r.Info.Notes
			}
		case "rollback":
			err = e.agents.helm.Rollback(step.App.Helm.Version, args)
		case "uninstall":
			var r *release.UninstallReleaseResponse
			r, err = e.agents.helm.Uninstall(args)
			output = r.Info
		case "upgrade":
			var r *release.Release
			r, err = e.agents.helm.Upgrade(args)
			if r != nil {
				output = r.Info.Notes
			}
		default:
			return Failed, append(l, log{
				level:   Error,
				message: "helm command not found",
			})
		}

		if err != nil {
			return Failed, append(l, log{
				level:   Error,
				message: err.Error(),
			})
		}

		l = append(l, log{
			level:   Info,
			message: logs.String(),
		})
		l = append(l, log{
			level:   Info,
			message: output,
		})

		return Succeeded, l
	}
	return Failed, append(l, log{
		level:   Error,
		message: "only helm apps are currently supported",
	})
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

func (e *RunEngine) executeHold(h string) error {
	d, err := time.ParseDuration(h)
	if err != nil {
		return err
	}
	time.Sleep(d)
	return nil
}

func (s Step) getType() string {
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
