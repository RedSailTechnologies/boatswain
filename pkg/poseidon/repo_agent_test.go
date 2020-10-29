package poseidon

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// NOTE AdamP - these are all lower t as they take a while to run due to
// downloads...I wonder if we can just mock?

func testCheckIndex(t *testing.T) {
	sut := defaultRepoAgent{}
	config := &RepoConfig{
		Name:     "stable",
		Endpoint: "https://charts.helm.sh/stable",
	}
	repo, _ := config.ToChartRepo()
	assert.True(t, sut.checkIndex(repo))
}

func testDownloadChart(t *testing.T) {
	config := &RepoConfig{
		Name:     "stable",
		Endpoint: "https://charts.helm.sh/stable",
	}

	name := "envoy"
	version := "1.0.0"
	opts := config.ToChartPathOptions()
	out := os.TempDir()
	defer os.Remove(out + "/envoy-1.0.0.tgz")

	sut := defaultRepoAgent{}
	file, err := sut.downloadChart(name, version, out, config.Endpoint, opts)
	assert.NotNil(t, file.Contents)
	assert.Nil(t, err)
}

func testGetCharts(t *testing.T) {
	sut := defaultRepoAgent{}
	config := &RepoConfig{
		Name:     "stable",
		Endpoint: "https://charts.helm.sh/stable",
	}
	repo, _ := config.ToChartRepo()
	out, err := sut.getCharts(repo)
	assert.NotNil(t, out)
	assert.Nil(t, err)
}
