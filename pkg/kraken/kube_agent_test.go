package kraken

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetClusterStatus(t *testing.T) {
	var fakeClientset kubernetes.Interface
	fakeClientset = fake.NewSimpleClientset(
		&v1.NodeList{
			Items: []v1.Node{
				v1.Node{
					Status: v1.NodeStatus{
						Conditions: []v1.NodeCondition{
							v1.NodeCondition{
								Type:   v1.NodeReady,
								Status: v1.ConditionTrue,
							},
						},
					},
				},
			},
		},
	)

	sut := &defaultKubeAgent{}
	assert.True(t, sut.getClusterStatus(fakeClientset, "doesntmatter"))
}
