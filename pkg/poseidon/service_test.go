package poseidon

import (
	"context"
	"errors"
	"os"
	"testing"

	pb "github.com/redsailtechnologies/boatswain/rpc/poseidon"
	"github.com/stretchr/testify/assert"
	"github.com/twitchtv/twirp"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/repo"
)

type mockRepoAgent struct{}

func (m mockRepoAgent) checkIndex(r *repo.ChartRepository) bool {
	if r.Config.URL == configA.Endpoint {
		return true
	}
	return false
}

func (m mockRepoAgent) downloadChart(n string, v string, d string, e string, o *action.ChartPathOptions) (*pb.File, error) {
	if n == "good-chart" {
		return &pb.File{
			Name:     "good-chart-0.1.0.tgz",
			Contents: []byte("some filey stuff"),
		}, nil
	}
	return nil, errors.New("some error here")
}

func (m mockRepoAgent) getCharts(*repo.ChartRepository) (map[string]repo.ChartVersions, error) {
	return nil, nil
}

var configA = RepoConfig{
	Name:     "repoA",
	Endpoint: "http://endpointa.com",
}

var configB = RepoConfig{
	Name:     "Brepo",
	Endpoint: "https://b-endpoint.com",
}

func TestNewCreatesDirectory(t *testing.T) {
	dir := os.TempDir() + "/poseidon-cache-dir-test"
	defer os.Remove(dir)
	New(&Config{
		Repos: []RepoConfig{
			configA,
			configB,
		},
		CacheDir: dir,
	})
}

func TestChartsErrorsOnNonexistentRepo(t *testing.T) {
	sut := New(&Config{
		Repos: []RepoConfig{
			RepoConfig{
				Name:     "bad",
				Endpoint: "repo.com",
			},
		},
		CacheDir: os.TempDir(),
	})

	resp, err := sut.Charts(context.TODO(), &pb.Repo{Name: "alsobad"})
	assert.Nil(t, resp)
	assert.Equal(t, twirp.InternalError("error getting repo config"), err)
}

func TestChartsErrorsOnBadRepo(t *testing.T) {
	sut := New(&Config{
		Repos: []RepoConfig{
			RepoConfig{
				Name:     "bad",
				Endpoint: "repo.com",
			},
		},
		CacheDir: os.TempDir(),
	})

	resp, err := sut.Charts(context.TODO(), &pb.Repo{Name: "bad"})
	assert.Nil(t, resp)
	assert.Equal(t, twirp.InternalError("error getting helm repo"), err)
}

func TestDownloadChartErrorsOnBadConfig(t *testing.T) {
	sut := New(&Config{
		Repos: []RepoConfig{
			RepoConfig{
				Name:     "bad",
				Endpoint: "repo.com",
			},
		},
		CacheDir: os.TempDir(),
	})

	resp, err := sut.DownloadChart(context.TODO(), &pb.DownloadRequest{})
	assert.Nil(t, resp)
	assert.Equal(t, twirp.InternalError("error getting repo config"), err)
}

func TestDownloadChartReportsAgentErrors(t *testing.T) {
	sut := New(&Config{
		Repos: []RepoConfig{
			configA,
		},
		CacheDir: os.TempDir(),
	})
	sut.repoAgent = mockRepoAgent{}

	a, aErr := sut.DownloadChart(context.TODO(), &pb.DownloadRequest{
		ChartName:    "good-chart",
		ChartVersion: "0.1.0",
		RepoName:     "repoA",
	})
	b, bErr := sut.DownloadChart(context.TODO(), &pb.DownloadRequest{
		ChartName:    "bad-chart",
		ChartVersion: "bad.0-.1",
		RepoName:     "repoA",
	})

	assert.Equal(t, "good-chart-0.1.0.tgz", a.Name)
	assert.Equal(t, []byte("some filey stuff"), a.Contents)
	assert.Nil(t, aErr)
	assert.Nil(t, b)
	assert.Equal(t, twirp.InternalError("error downloading chart"), bErr)
}

func TestReposCallsAgent(t *testing.T) {
	sut := New(&Config{
		Repos: []RepoConfig{
			configA,
			configB,
		},
		CacheDir: os.TempDir(),
	})
	sut.repoAgent = mockRepoAgent{}
	response, err := sut.Repos(context.TODO(), &pb.ReposRequest{})

	assert.Nil(t, err)
	assert.True(t, response.Repos[0].Ready)
	assert.False(t, response.Repos[1].Ready)
}

func TestReposErrorsOnBadConfig(t *testing.T) {
	sut := New(&Config{
		Repos: []RepoConfig{
			RepoConfig{
				Name:     "bad",
				Endpoint: "repo.com",
			},
		},
		CacheDir: os.TempDir(),
	})

	resp, err := sut.Repos(context.TODO(), &pb.ReposRequest{})
	assert.Nil(t, resp)
	assert.Equal(t, twirp.InternalError("error getting helm repo"), err)
}

func TestBuildChartURL(t *testing.T) {
	a := buildChartURL("https://somechartrepo.com", "chart-0.1.0.tgz")
	b := buildChartURL("http://anotherchartrepo.com/", "chart-0.1.0.tgz")
	assert.Equal(t, "https://somechartrepo.com/chart-0.1.0.tgz", a)
	assert.Equal(t, "http://anotherchartrepo.com/chart-0.1.0.tgz", b)
}
