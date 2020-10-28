package kraken

import (
	"bytes"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/release"

	"github.com/redsailtechnologies/boatswain/rpc/poseidon"
)

type helmAgent interface {
	getReleases(*action.Configuration, string) ([]*release.Release, error)
	upgradeRelease(*action.Configuration, string, *poseidon.File, string, map[string]interface{}) (*release.Release, error)
}

type defaultHelmAgent struct{}

func (h defaultHelmAgent) getReleases(cfg *action.Configuration, cluster string) ([]*release.Release, error) {
	list := action.NewList(cfg)
	list.All = true
	list.AllNamespaces = true
	list.Limit = 0
	list.SetStateMask()

	releases, err := list.Run()
	if err != nil {
		return nil, err
	}

	return releases, nil
}

func (h defaultHelmAgent) upgradeRelease(cfg *action.Configuration, n string, f *poseidon.File, ns string, vals map[string]interface{}) (*release.Release, error) {
	chart, err := loader.LoadArchive(bytes.NewReader(f.Contents))
	if err != nil {
		return nil, err
	}

	upgrade := action.NewUpgrade(cfg)
	upgrade.Namespace = ns
	return upgrade.Run(n, chart, vals)
}
