package kube

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/redsailtechnologies/boatswain/pkg/logger"
)

// Agent is used for communication with Kubernetes, made into an interface for testability
type Agent interface {
	GetDeployments(args *Args) (*Result, error)
	GetStatefulSets(args *Args) (*Result, error)
	GetStatus(args *Args) (*Result, error)
}

// AgentAction is an enum/alias to make calling methods typed
type AgentAction string

const (
	// GetDeployments represents the GetDeployments method
	GetDeployments AgentAction = "GetDeployments"

	// GetStatefulSets represents the GetStatefulSets method
	GetStatefulSets AgentAction = "GetStatefulSets"

	// GetStatus represents the GetStatus method
	GetStatus AgentAction = "GetStatus"
)

// DefaultAgent is the default implementation of the KubeAgent
type DefaultAgent struct {
	kube kubernetes.Interface
}

// NewDefaultAgent inits the default agent with the specified kube interface
func NewDefaultAgent(kube kubernetes.Interface) *DefaultAgent {
	return &DefaultAgent{
		kube: kube,
	}
}

// GetDeployments gets all the deployments for a particular cluster
func (k DefaultAgent) GetDeployments(args *Args) (*Result, error) {
	d, err := k.kube.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("could not get deployments from cluster", "error", err)
		return nil, err
	}
	return &Result{
		Data: d.Items,
		Type: DeploymentsResult,
	}, nil
}

// GetStatefulSets gets all the statefulsets for a particular cluster
func (k DefaultAgent) GetStatefulSets(args *Args) (*Result, error) {
	ss, err := k.kube.AppsV1().StatefulSets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("could not get statefulsets from cluster", "error", err)
		return nil, err
	}
	return &Result{
		Data: ss.Items,
		Type: StatefulSetsResult,
	}, nil
}

// GetStatus returns the status of a cluster by ensuring each node is in a ready state
func (k DefaultAgent) GetStatus(args *Args) (*Result, error) {
	nodes, err := k.kube.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("could not get nodes from cluster", "error", err)
		return nil, err
	}

	for _, node := range nodes.Items {
		for _, condition := range node.Status.Conditions {
			if condition.Type == "Ready" {
				if condition.Status != "True" {
					return &Result{
						Data: false,
						Type: StatusResult,
					}, nil
				}
			}
		}
	}

	return &Result{
		Data: true,
		Type: StatusResult,
	}, nil
}
