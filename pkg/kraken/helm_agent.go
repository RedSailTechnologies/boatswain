package kraken

import (
	"errors"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

type helmAgent interface {
	getReleases(*action.Configuration, string) ([]*release.Release, error)
	getReleaseStatus(*action.Configuration, string) (*release.Release, error)
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
		return nil, errors.New("could not list releases")
	}

	return releases, nil
}

func (h defaultHelmAgent) getReleaseStatus(cfg *action.Configuration, cluster string) (*release.Release, error) {
	return nil, errors.New("not implemented")
}
