package run

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/redsailtechnologies/boatswain/pkg/cluster"
	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/deployment/template"
	"github.com/redsailtechnologies/boatswain/pkg/helm"
	"github.com/redsailtechnologies/boatswain/pkg/kube"
	"github.com/redsailtechnologies/boatswain/pkg/repo"
	"github.com/redsailtechnologies/boatswain/rpc/agent"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"
)

func (e *Engine) executeHelmStep(step *template.Step) Status {
	var err error
	clusters, err := e.clusters.All()
	if err != nil {
		e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
		return Failed
	}
	cl := getCluster(step.Helm.Cluster, clusters)
	if cl == nil {
		e.run.AppendLog("cluster not found", Error, ddd.NewTimestamp())
		return Failed
	}

	var chart []byte
	if step.Helm.Command == "install" || step.Helm.Command == "upgrade" {
		chart, err = e.downloadChart(step.Helm.Chart.Name, step.Helm.Chart.Version, step.Helm.Chart.Repo)
		if err != nil {
			e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
			return Failed
		}
	}

	var vals map[string]interface{}
	if step.Helm.Values != nil {
		if step.Helm.Values.Library != nil {
			lib := *step.Helm.Values.Library
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

		if step.Helm.Values.Raw != nil {
			// order is important - favor raw values
			vals = mergeVals(*step.Helm.Values.Raw, vals)
		}
	}

	name, err := e.getReleaseName(step, cl)
	if err != nil {
		e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
		return Failed
	}

	args := helm.Args{
		Name:      name,
		Namespace: step.Helm.Namespace,
		Chart:     chart,
		Values:    vals,
		Wait:      step.Helm.Options.Wait,
		Timeout:   step.Helm.Options.Timeout,
		Install:   step.Helm.Options.Install,
		Version:   step.Helm.Options.RollbackVersion,
	}
	jsonArgs, err := json.Marshal(args)
	if err != nil {
		e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
		return Failed
	}

	action := &agent.Action{
		Uuid:           ddd.NewUUID(),
		ClusterUuid:    cl.UUID(),
		ClusterToken:   cl.Token(),
		ActionType:     agent.ActionType_HELM_ACTION,
		TimeoutSeconds: 600, // FIXME - configurable?
		Args:           jsonArgs,
	}

	var result *agent.Result
	switch step.Helm.Command {
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

func (e *Engine) downloadChart(c, v, r string) ([]byte, error) {
	if c == "" || v == "" || r == "" {
		return nil, errors.New("cannot download chart without chart, version, and name")
	}
	allRepos, err := e.repos.All()
	if err != nil {
		return nil, err
	}

	repo, err := findRepo(r, allRepos)
	if err != nil {
		return nil, err
	}

	if repo.HelmOCI() {
		return e.repo.GetChartFromOCIV2(c, v, repo.Endpoint(), repo.Username(), repo.Password())
	}
	return e.repo.GetChart(c, v, repo.Endpoint(), repo.Token())
}

func (e *Engine) getReleaseName(s *template.Step, c *cluster.Cluster) (string, error) {
	name := s.Helm.Name
	if name == "" {
		if s.Helm.Selector == nil {
			return "", errors.New("no name or selector found")
		}

		jsonArgs, err := json.Marshal(kube.Args{
			Labels:    s.Helm.Selector,
			Namespace: s.Helm.Namespace,
		})
		if err != nil {
			return "", err
		}

		res, err := e.agent.Run(context.Background(), &agent.Action{
			Uuid:           ddd.NewUUID(),
			ClusterUuid:    c.UUID(),
			ClusterToken:   c.Token(),
			ActionType:     agent.ActionType_KUBE_ACTION,
			Action:         string(kube.GetReleaseName),
			TimeoutSeconds: 600, // FIXME - configurable?
			Args:           jsonArgs,
		})
		if err != nil {
			return "", err
		}

		name, err = kube.ConvertReleaseName(res.Data)
		if err != nil {
			return "", err
		}
	}
	return name, nil
}

func findRepo(r string, l []*repo.Repo) (*repo.Repo, error) {
	for _, repo := range l {
		if repo.Name() == r {
			return repo, nil
		}
	}
	return nil, errors.New("repo not found")
}

func getCluster(c string, l []*cluster.Cluster) *cluster.Cluster {
	for _, cluster := range l {
		if cluster.Name() == c {
			return cluster
		}
	}
	return nil
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
