package poseidon

import (
	"helm.sh/helm/v3/pkg/repo"
)

type repoAgent interface {
	checkIndex(*repo.ChartRepository) bool
}

type defaultRepoAgent struct{}

func (a defaultRepoAgent) checkIndex(r *repo.ChartRepository) bool {
	str, err := r.DownloadIndexFile()
	return str != "" && err == nil
}
