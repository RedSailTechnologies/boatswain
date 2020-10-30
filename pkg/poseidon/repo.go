package poseidon

import (
	pb "github.com/redsailtechnologies/boatswain/rpc/poseidon"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
)

// Repo is a type def to include the protobuf definition so we can extend it
type Repo struct {
	*pb.Repo
}

// ToChartPathOptions returns the chart path options for helm
func (r *Repo) ToChartPathOptions() *action.ChartPathOptions {
	return &action.ChartPathOptions{
		InsecureSkipTLSverify: true,
		RepoURL:               r.Endpoint,
	}
}

// ToChartRepo takes the repo and makes it into the helm repo struct
func (r *Repo) ToChartRepo() (*repo.ChartRepository, error) {
	providers := []getter.Provider{
		getter.Provider{
			Schemes: []string{"http", "https"},
			New:     getter.NewHTTPGetter,
		},
	}

	entry := &repo.Entry{
		Name: r.Name,
		URL:  r.Endpoint,
		// TODO AdamP - we definitely want to support this soon!
		InsecureSkipTLSverify: true,
	}

	return repo.NewChartRepository(entry, providers)
}
