package poseidon

import (
	"testing"

	pb "github.com/redsailtechnologies/boatswain/rpc/poseidon"
	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/repo"
)

func TestToChartPathOptionsUsesCorrectEndpoint(t *testing.T) {
	endpoint := "http://ahelmrepo.notreal/"
	sut := &Repo{
		&pb.Repo{
			Name:     "helm-repo",
			Endpoint: endpoint,
		},
	}
	cpo := sut.ToChartPathOptions()
	assert.True(t, cpo.InsecureSkipTLSverify)
	assert.Equal(t, endpoint, cpo.RepoURL)
}

func TestToChartRepoUsesConfig(t *testing.T) {
	name := "helm-repo"
	endpoint := "http://ahelmrepo.notreal/"
	sut := &Repo{
		&pb.Repo{
			Name:     name,
			Endpoint: endpoint,
		},
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
	a := &Repo{
		&pb.Repo{
			Name:     "repoA",
			Endpoint: "endpoint.com",
		},
	}
	b := &Repo{
		&pb.Repo{
			Name:     "RepoB",
			Endpoint: "anotherendpoint.com",
		},
	}
	sut := &Service{
		repos: []*Repo{
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
