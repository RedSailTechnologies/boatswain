package helm

import (
	"helm.sh/helm/v3/pkg/repo"
)

// Agent is the interface we use to talk to helm packages
type Agent interface {
	CheckIndex(*repo.ChartRepository) bool
	GetCharts(*repo.ChartRepository) (map[string]repo.ChartVersions, error)
}

// DefaultAgent is the default implementation of the Agent interface
type DefaultAgent struct{}

// CheckIndex checks the index.yaml file at the repo's endpoint
func (a DefaultAgent) CheckIndex(r *repo.ChartRepository) bool {
	str, err := r.DownloadIndexFile()
	return str != "" && err == nil
}

// GetCharts gets all charts from a particular chart repo
func (a DefaultAgent) GetCharts(r *repo.ChartRepository) (map[string]repo.ChartVersions, error) {
	str, err := r.DownloadIndexFile()
	if err != nil {
		return nil, err
	}

	idx, err := repo.LoadIndexFile(str)
	if err != nil {
		return nil, err
	}

	return idx.Entries, nil
}

func getFullName(n string, v string) string {
	return n + "-" + v + ".tgz"
}
