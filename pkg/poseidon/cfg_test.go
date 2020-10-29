package poseidon

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/repo"
)

func TestToChartPathOptionsUsesCorrectEndpoint(t *testing.T) {
	endpoint := "http://ahelmrepo.notreal/"
	sut := &RepoConfig{
		Name:     "helm-repo",
		Endpoint: endpoint,
	}
	cpo := sut.ToChartPathOptions()
	assert.True(t, cpo.InsecureSkipTLSverify)
	assert.Equal(t, endpoint, cpo.RepoURL)
}

func TestToChartRepoUsesConfig(t *testing.T) {
	name := "helm-repo"
	endpoint := "http://ahelmrepo.notreal/"
	sut := &RepoConfig{
		Name:     name,
		Endpoint: endpoint,
	}
	out, err := sut.ToChartRepo()
	if err != nil {
		t.Error("error converting config to chart repo")
	}

	assert.Equal(t, &repo.Entry{
		Name:                  name,
		URL:                   endpoint,
		InsecureSkipTLSverify: true,
	}, out.Config)
}

func TestGetRepoConfig(t *testing.T) {
	a := RepoConfig{
		Name:     "repoA",
		Endpoint: "endpoint.com",
	}
	b := RepoConfig{
		Name:     "RepoB",
		Endpoint: "anotherendpoint.com",
	}
	sut := &Config{
		Repos: []RepoConfig{
			a,
			b,
		},
	}

	aOut, aErr := sut.getRepoConfig("repoA")
	bOut, bErr := sut.getRepoConfig("RepoB")
	cOut, cErr := sut.getRepoConfig("errorExpected")

	assert.Equal(t, "endpoint.com", aOut.Endpoint)
	assert.Nil(t, aErr)
	assert.Equal(t, "anotherendpoint.com", bOut.Endpoint)
	assert.Nil(t, bErr)
	assert.Nil(t, cOut)
	assert.NotNil(t, cErr)
}
