package kube

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetStatus(t *testing.T) {
	fakeClientset := func() (kubernetes.Interface, error) {
		return fake.NewSimpleClientset(
			&v1.NodeList{
				Items: []v1.Node{
					{
						Status: v1.NodeStatus{
							Conditions: []v1.NodeCondition{
								{
									Type:   v1.NodeReady,
									Status: v1.ConditionTrue,
								},
							},
						},
					},
				},
			},
		), nil
	}

	sut := NewDefaultAgent(fakeClientset)
	result, err := sut.GetStatus(&Args{})
	assert.True(t, result.Data.(bool))
	assert.Nil(t, err)
}
