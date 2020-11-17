package poseidon

import (
	"os"
	"testing"

	pb "github.com/redsailtechnologies/boatswain/rpc/poseidon"
	"github.com/stretchr/testify/assert"
)

// NOTE AdamP - these are all lower t as they take a while to run due to
// downloads...I wonder if we can just mock?

func testCheckIndex(t *testing.T) {
	sut := defaultRepoAgent{}
	repo := &Repo{
		&pb.Repo{
			Name:     "stable",
			Endpoint: "https://charts.helm.sh/stable",
		},
	}
	chartRepo, _ := repo.ToChartRepo()
	assert.True(t, sut.checkIndex(chartRepo))
}

func testDownloadChart(t *testing.T) {
	repo := &Repo{
		&pb.Repo{
			Name:     "stable",
			Endpoint: "https://charts.helm.sh/stable",
		},
	}

	name := "envoy"
	version := "1.0.0"
	opts := repo.ToChartPathOptions()
	out := os.TempDir()
	defer os.Remove(out + "/envoy-1.0.0.tgz")

	sut := defaultRepoAgent{}
	file, err := sut.downloadChart(name, version, out, repo.Endpoint, opts)
	assert.NotNil(t, file.Contents)
	assert.Nil(t, err)
}

func testGetCharts(t *testing.T) {
	sut := defaultRepoAgent{}
	repo := &Repo{
		&pb.Repo{
			Name:     "stable",
			Endpoint: "https://charts.helm.sh/stable",
		},
	}
	chartRepo, _ := repo.ToChartRepo()
	out, err := sut.getCharts(chartRepo)
	assert.NotNil(t, out)
	assert.Nil(t, err)
}
