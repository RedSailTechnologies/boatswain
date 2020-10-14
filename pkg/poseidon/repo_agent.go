package poseidon

import (
	"helm.sh/helm/v3/pkg/repo"
)

type repoAgent interface {
	checkIndex(*repo.ChartRepository) bool
	getCharts(*repo.ChartRepository) (map[string]repo.ChartVersions, error)
}

type defaultRepoAgent struct{}

func (a defaultRepoAgent) checkIndex(r *repo.ChartRepository) bool {
	str, err := r.DownloadIndexFile()
	return str != "" && err == nil
}

func (a defaultRepoAgent) getCharts(r *repo.ChartRepository) (map[string]repo.ChartVersions, error) {
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
