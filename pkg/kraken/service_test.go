package kraken

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/client-go/kubernetes"

	pb "github.com/redsailtechnologies/boatswain/rpc/kraken"
)

type MockedKubeAgent struct {
	mock.Mock
}

func (m *MockedKubeAgent) getClusterStatus(kube kubernetes.Interface, name string) bool {
	args := m.Called(kube, name)
	return args.Bool(0)
}

func Test_ClustersGetsAll(t *testing.T) {
	firstConfig := ClusterConfig{
		Name:     "first",
		Endpoint: "someendpoint",
		Token:    "tokenA",
		Cert:     "notreal",
	}
	secondConfig := ClusterConfig{
		Name:     "second",
		Endpoint: "anotherendpoint",
		Token:    "tokenB",
		Cert:     "alsonotreal",
	}

	mockAgent := MockedKubeAgent{}
	mockAgent.On("getClusterStatus", mock.Anything, mock.Anything).Return(true)

	sut := &Service{
		config: &Config{
			Clusters: []ClusterConfig{
				firstConfig,
				secondConfig,
			},
		},
		agent: &mockAgent,
	}

	result, err := sut.Clusters(nil, &pb.ClustersRequest{})
	assert.Equal(t, 2, len(result.Clusters))
	assert.True(t, result.Clusters[0].Ready)
	assert.True(t, result.Clusters[1].Ready)
	assert.Nil(t, err)
}

func Test_ClustersReportsCorrectStatus(t *testing.T) {
	firstConfig := ClusterConfig{
		Name:     "first",
		Endpoint: "someendpoint",
		Token:    "tokenA",
		Cert:     "notreal",
	}
	secondConfig := ClusterConfig{
		Name:     "second",
		Endpoint: "anotherendpoint",
		Token:    "tokenB",
		Cert:     "alsonotreal",
	}

	mockAgent := MockedKubeAgent{}
	mockAgent.On("getClusterStatus", mock.Anything, "first").Return(true)
	mockAgent.On("getClusterStatus", mock.Anything, "second").Return(false)

	sut := &Service{
		config: &Config{
			Clusters: []ClusterConfig{
				firstConfig,
				secondConfig,
			},
		},
		agent: &mockAgent,
	}

	result, err := sut.Clusters(nil, &pb.ClustersRequest{})
	assert.Equal(t, 2, len(result.Clusters))
	assert.True(t, result.Clusters[0].Ready)
	assert.False(t, result.Clusters[1].Ready)
	assert.Nil(t, err)
}
