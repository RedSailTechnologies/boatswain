package kraken

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/client-go/kubernetes"

	"github.com/redsailtechnologies/boatswain/pkg/poseidon"
	pb "github.com/redsailtechnologies/boatswain/rpc/kraken"
	poseidonPB "github.com/redsailtechnologies/boatswain/rpc/poseidon"
)

type mockedKubeAgent struct {
	mock.Mock
}

func (m *mockedKubeAgent) getClusterStatus(kube kubernetes.Interface, name string) bool {
	args := m.Called(kube, name)
	return args.Bool(0)
}

type mockedHelmAgent struct {
	mock.Mock
}

func (m *mockedHelmAgent) getReleases(cfg *action.Configuration, cluster string) ([]*release.Release, error) {
	args := m.Called(cfg, cluster)
	return args.Get(0).([]*release.Release), nil
}

func (m *mockedHelmAgent) upgradeRelease(cfg *action.Configuration, n string, f *poseidonPB.File, ns string, vals map[string]interface{}) (*release.Release, error) {
	return nil, errors.New("not implemented")
}

var firstConfig = ClusterConfig{
	Name:     "first",
	Endpoint: "someendpoint",
	Token:    "tokenA",
	Cert:     "notreal",
}
var secondConfig = ClusterConfig{
	Name:     "second",
	Endpoint: "anotherendpoint",
	Token:    "tokenB",
	Cert:     "alsonotreal",
}

var firstCluster = pb.Cluster{
	Name:     firstConfig.Name,
	Endpoint: firstConfig.Endpoint,
	Ready:    false,
}
var secondCluster = pb.Cluster{
	Name:     secondConfig.Name,
	Endpoint: secondConfig.Endpoint,
	Ready:    true,
}

func TestNew(t *testing.T) {
	sut := New(&Config{
		Clusters: []ClusterConfig{
			ClusterConfig{
				Name:     "first",
				Endpoint: "someendpoint",
				Token:    "tokenA",
				Cert:     "notreal",
			},
		},
	}, &poseidon.Service{})
	assert.NotNil(t, sut)
}

func TestClustersGetsAll(t *testing.T) {
	mockAgent := mockedKubeAgent{}
	mockAgent.On("getClusterStatus", mock.Anything, mock.Anything).Return(true)

	sut := &Service{
		config: &Config{
			Clusters: []ClusterConfig{
				firstConfig,
				secondConfig,
			},
		},
		kubeAgent: &mockAgent,
	}

	result, err := sut.Clusters(nil, &pb.ClustersRequest{})
	assert.Equal(t, 2, len(result.Clusters))
	assert.True(t, result.Clusters[0].Ready)
	assert.True(t, result.Clusters[1].Ready)
	assert.Nil(t, err)
}

func TestClustersReportsCorrectStatus(t *testing.T) {
	mockAgent := mockedKubeAgent{}
	mockAgent.On("getClusterStatus", mock.Anything, "first").Return(true)
	mockAgent.On("getClusterStatus", mock.Anything, "second").Return(false)

	sut := &Service{
		config: &Config{
			Clusters: []ClusterConfig{
				firstConfig,
				secondConfig,
			},
		},
		kubeAgent: &mockAgent,
	}

	result, err := sut.Clusters(nil, &pb.ClustersRequest{})
	assert.Equal(t, 2, len(result.Clusters))
	assert.True(t, result.Clusters[0].Ready)
	assert.False(t, result.Clusters[1].Ready)
	assert.Nil(t, err)
}

func TestClusterStatusUpdatesClusterCorrectly(t *testing.T) {
	mockAgent := mockedKubeAgent{}
	mockAgent.On("getClusterStatus", mock.Anything, firstConfig.Name).Return(true)
	mockAgent.On("getClusterStatus", mock.Anything, secondConfig.Name).Return(false)

	sut := &Service{
		config: &Config{
			Clusters: []ClusterConfig{
				firstConfig,
				secondConfig,
			},
		},
		kubeAgent: &mockAgent,
	}

	first, err := sut.ClusterStatus(nil, &firstCluster)
	assert.True(t, first.Ready)
	assert.Nil(t, err)

	second, err := sut.ClusterStatus(nil, &secondCluster)
	assert.False(t, second.Ready)
	assert.Nil(t, err)
}

func TestReleasesSortsReleases(t *testing.T) {
	firstRelease := &release.Release{
		Name:      "firstRelease",
		Namespace: "first-namespace",
		Info: &release.Info{
			Status: release.StatusUnknown,
		},
		Chart: &chart.Chart{
			Metadata: &chart.Metadata{
				Name:       "first-chart",
				Version:    "1",
				AppVersion: "2",
			},
		},
	}
	secondRelease := &release.Release{
		Name:      "secondRelease",
		Namespace: "second-namespace",
		Info: &release.Info{
			Status: release.StatusDeployed,
		},
		Chart: &chart.Chart{
			Metadata: &chart.Metadata{
				Name:       "second-chart",
				Version:    "3",
				AppVersion: "4",
			},
		},
	}
	thirdRelease := &release.Release{
		Name:      "secondRelease",
		Namespace: "third-namespace",
		Info: &release.Info{
			Status: release.StatusPendingInstall,
		},
		Chart: &chart.Chart{
			Metadata: &chart.Metadata{
				Name:       "second-chart",
				Version:    "5",
				AppVersion: "6",
			},
		},
	}

	mockAgent := mockedHelmAgent{}
	sut := &Service{
		config: &Config{
			Clusters: []ClusterConfig{
				firstConfig,
				secondConfig,
			},
		},
		helmAgent: &mockAgent,
	}

	// firstClient, _ := sut.config.ToHelmClient(firstConfig.Name)
	// secondClient, _ := sut.config.ToHelmClient(secondConfig.Name)

	mockAgent.On("getReleases", mock.Anything, firstConfig.Name).Return([]*release.Release{
		firstRelease,
		secondRelease,
	}, nil)
	mockAgent.On("getReleases", mock.Anything, secondConfig.Name).Return([]*release.Release{
		thirdRelease,
	}, nil)

	response, err := sut.Releases(nil, &pb.ReleaseRequest{
		Clusters: []*pb.Cluster{
			&firstCluster,
			&secondCluster,
		},
	})

	assert.Nil(t, err)
	assert.Equal(t, 2, len(response.ReleaseLists))
	assert.Equal(t, 1, len(response.ReleaseLists[0].Releases))
	assert.Equal(t, 2, len(response.ReleaseLists[1].Releases))
	assert.Equal(t, firstRelease.Name, response.ReleaseLists[0].Name)
	assert.Equal(t, secondRelease.Name, response.ReleaseLists[1].Name)
	assert.Equal(t, firstRelease.Chart.Metadata.Version, response.ReleaseLists[0].Releases[0].ChartVersion)
	assert.Equal(t, secondRelease.Namespace, response.ReleaseLists[1].Releases[0].Namespace)
	assert.NotEqual(t, secondRelease.Namespace, response.ReleaseLists[1].Releases[1].Namespace)
	assert.Equal(t, thirdRelease.Chart.Metadata.AppVersion, response.ReleaseLists[1].Releases[1].AppVersion)
}
