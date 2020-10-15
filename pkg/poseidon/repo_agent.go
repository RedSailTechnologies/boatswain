package poseidon

import (
	"io/ioutil"
	"os"
	"path"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/repo"

	pb "github.com/redsailtechnologies/boatswain/rpc/poseidon"
)

type repoAgent interface {
	checkIndex(*repo.ChartRepository) bool
	downloadChart(string, string, string, string, *action.ChartPathOptions) (*pb.File, error)
	getCharts(*repo.ChartRepository) (map[string]repo.ChartVersions, error)
}

type defaultRepoAgent struct{}

func (a defaultRepoAgent) checkIndex(r *repo.ChartRepository) bool {
	str, err := r.DownloadIndexFile()
	return str != "" && err == nil
}

func (a defaultRepoAgent) downloadChart(name, version, out, endpoint string, opts *action.ChartPathOptions) (*pb.File, error) {
	pull := action.NewPull()
	pull.ChartPathOptions = *opts
	pull.Settings = cli.New()
	pull.RepoURL = endpoint
	pull.Version = version

	// TODO AdamP - we may want to implement some kind of directory management to prevent overflow here..
	pull.DestDir = out

	if _, err := pull.Run(name); err != nil {
		return nil, err
	}

	file, err := os.Open(path.Join(out, getFullName(name, version)))
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return &pb.File{
		Name:     name,
		Contents: bytes,
	}, nil
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

func getFullName(n string, v string) string {
	return n + "-" + v + ".tgz"
}
