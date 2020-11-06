package kraken

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/redsailtechnologies/boatswain/pkg/poseidon"
	pb "github.com/redsailtechnologies/boatswain/rpc/kraken"
	poseidonPB "github.com/redsailtechnologies/boatswain/rpc/poseidon"
)

type mockedKubeAgent struct {
	mock.Mock
}

func (m *mockedKubeAgent) getClusterDeployments(kubernetes.Interface, string) ([]appsv1.Deployment, error) {
	return nil, nil
}

func (m *mockedKubeAgent) getClusterStatefulSets(kubernetes.Interface, string) ([]appsv1.StatefulSet, error) {
	return nil, nil
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

var firstCluster = &Cluster{
	&pb.Cluster{
		Name:     "first",
		Endpoint: "someendpoint",
		Token:    "tokenA",
		Cert:     "notreal",
		Ready:    false,
	},
}

var secondCluster = &Cluster{
	&pb.Cluster{
		Name:     "second",
		Endpoint: "anotherendpoint",
		Token:    "tokenB",
		Cert:     "alsonotreal",
		Ready:    true,
	},
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

func TestAddEditDeleteArraySizing(t *testing.T) {
	sut := &Service{}

	sut.AddCluster(context.TODO(), secondCluster.Cluster)
	sut.AddCluster(context.TODO(), secondCluster.Cluster)
	assert.Len(t, sut.clusters, 2)

	sut.DeleteCluster(context.TODO(), sut.clusters[0].Cluster)
	assert.Len(t, sut.clusters, 1)

	sut.EditCluster(context.TODO(), sut.clusters[0].Cluster)
	assert.Len(t, sut.clusters, 1)

	sut.DeleteCluster(context.TODO(), sut.clusters[0].Cluster)
	assert.Len(t, sut.clusters, 0)
}

func TestAddOnlyAddsCorrectFields(t *testing.T) {
	sut := &Service{}
	sut.AddCluster(context.TODO(), secondCluster.Cluster)

	assert.Equal(t, "second", sut.clusters[0].Name)
	assert.Equal(t, "anotherendpoint", sut.clusters[0].Endpoint)
	assert.Equal(t, "tokenB", sut.clusters[0].Token)
	assert.Equal(t, "alsonotreal", sut.clusters[0].Cert)
	assert.NotEqual(t, secondCluster.Uuid, sut.clusters[0].Uuid)
	assert.NotEqual(t, secondCluster.Ready, sut.clusters[0].Ready)
}

func TestEditCorrectlyModifies(t *testing.T) {
	sut := &Service{}
	sut.AddCluster(context.TODO(), secondCluster.Cluster)

	orig := sut.clusters[0]
	copy := &Cluster{
		&pb.Cluster{
			Uuid:     orig.Uuid,
			Name:     orig.Name,
			Endpoint: orig.Endpoint,
			Token:    orig.Token,
			Cert:     orig.Cert,
			Ready:    orig.Ready,
		},
	}
	copy.Name = "newname"

	assert.NotEqual(t, copy.Name, sut.clusters[0].Name)
	sut.EditCluster(context.TODO(), copy.Cluster)

	assert.Equal(t, "newname", sut.clusters[0].Name)
}

func TestDeleteRemovesCorrectCluster(t *testing.T) {
	sut := &Service{}
	sut.AddCluster(context.TODO(), secondCluster.Cluster)
	sut.AddCluster(context.TODO(), firstCluster.Cluster)
	sut.AddCluster(context.TODO(), secondCluster.Cluster)

	delete := *sut.clusters[2]
	assert.NotEqual(t, delete.Uuid, sut.clusters[0].Uuid)
	sut.DeleteCluster(context.TODO(), delete.Cluster)

	assert.Len(t, sut.clusters, 2)
	assert.NotEqual(t, delete.Uuid, sut.clusters[0].Uuid)
	assert.NotEqual(t, delete.Uuid, sut.clusters[1].Uuid)
}

func TestClustersGetsAll(t *testing.T) {
	mockAgent := mockedKubeAgent{}
	mockAgent.On("getClusterStatus", mock.Anything, mock.Anything).Return(true)

	sut := &Service{
		clusters: []*Cluster{
			firstCluster,
			secondCluster,
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
		clusters: []*Cluster{
			firstCluster,
			secondCluster,
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
	mockAgent.On("getClusterStatus", mock.Anything, firstCluster.Name).Return(true)
	mockAgent.On("getClusterStatus", mock.Anything, secondCluster.Name).Return(false)

	sut := &Service{
		clusters: []*Cluster{
			firstCluster,
			secondCluster,
		},
		kubeAgent: &mockAgent,
	}

	first, err := sut.ClusterStatus(nil, firstCluster.Cluster)
	assert.True(t, first.Ready)
	assert.Nil(t, err)

	second, err := sut.ClusterStatus(nil, secondCluster.Cluster)
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
		clusters: []*Cluster{
			firstCluster,
			secondCluster,
		},
		helmAgent: &mockAgent,
	}

	mockAgent.On("getReleases", mock.Anything, firstCluster.Name).Return([]*release.Release{
		firstRelease,
		secondRelease,
	}, nil)
	mockAgent.On("getReleases", mock.Anything, secondCluster.Name).Return([]*release.Release{
		thirdRelease,
	}, nil)

	response, err := sut.Releases(nil, &pb.ReleaseRequest{
		Clusters: []*pb.Cluster{
			firstCluster.Cluster,
			secondCluster.Cluster,
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

func TestGetClusterByName(t *testing.T) {
	cluster := &Cluster{
		&pb.Cluster{
			Name:     "cluster",
			Endpoint: "www.not.real",
			Token:    "abcdefg",
			Cert:     "notarealcert...",
		},
	}
	sut := Service{
		clusters: []*Cluster{
			cluster,
		},
	}

	valid, noErr := sut.getClusterByName("cluster")
	invalid, shouldErr := sut.getClusterByName("doesn'texist")

	assert.Equal(t, cluster, valid)
	assert.Nil(t, noErr)
	assert.Nil(t, invalid)
	assert.Equal(t, errors.New("cluster not found"), shouldErr)
}
