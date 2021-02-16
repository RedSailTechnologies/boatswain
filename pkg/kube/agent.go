package kube

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/redsailtechnologies/boatswain/pkg/logger"
)

// Agent is used for communication with Kubernetes, made into an interface for testability
type Agent interface {
	GetDeployments() ([]appsv1.Deployment, error)
	GetStatefulSets() ([]appsv1.StatefulSet, error)
	GetStatus() bool
}

// DefaultAgent is the default implementation of the KubeAgent
type DefaultAgent struct {
	kube kubernetes.Interface
}

// Args are arguments common to all commands
type Args struct {
	Name string
}

// NewDefaultAgent inits the default agent with the specified kube interface
func NewDefaultAgent(kube kubernetes.Interface) *DefaultAgent {
	return &DefaultAgent{
		kube: kube,
	}
}

// GetDeployments gets all the deployments for a particular cluster
func (k DefaultAgent) GetDeployments() ([]appsv1.Deployment, error) {
	d, err := k.kube.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("could not get deployments from cluster", "error", err)
		return nil, err
	}
	return d.Items, nil
}

// GetStatefulSets gets all the statefulsets for a particular cluster
func (k DefaultAgent) GetStatefulSets() ([]appsv1.StatefulSet, error) {
	ss, err := k.kube.AppsV1().StatefulSets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("could not get statefulsets from cluster", "error", err)
		return nil, err
	}
	return ss.Items, nil
}

// GetStatus returns the status of a cluster by ensuring each node is in a ready state
func (k DefaultAgent) GetStatus() bool {
	nodes, err := k.kube.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("could not get nodes from cluster", "error", err)
		return false
	}

	for _, node := range nodes.Items {
		for _, condition := range node.Status.Conditions {
			if condition.Type == "Ready" {
				if condition.Status != "True" {
					return false
				}
			}
		}
	}

	return true
}
