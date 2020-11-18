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
	GetClusterDeployments(kubernetes.Interface, string) ([]appsv1.Deployment, error)
	GetClusterStatefulSets(kubernetes.Interface, string) ([]appsv1.StatefulSet, error)
	GetClusterStatus(kubernetes.Interface, string) bool
}

// DefaultAgent is the default implementation of the KubeAgent
type DefaultAgent struct{}

// GetClusterDeployments gets all the deployments for a particular cluster
func (k DefaultAgent) GetClusterDeployments(kube kubernetes.Interface, name string) ([]appsv1.Deployment, error) {
	d, err := kube.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("could not get deployments from cluster", "cluster", name, "error", err)
		return nil, err
	}
	return d.Items, nil
}

// GetClusterStatefulSets gets all the statefulsets for a particular cluster
func (k DefaultAgent) GetClusterStatefulSets(kube kubernetes.Interface, name string) ([]appsv1.StatefulSet, error) {
	ss, err := kube.AppsV1().StatefulSets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("could not get statefulsets from cluster", "cluster", name, "error", err)
		return nil, err
	}
	return ss.Items, nil
}

// GetClusterStatus returns the status of a cluster by ensuring each node is in a ready state
func (k DefaultAgent) GetClusterStatus(kube kubernetes.Interface, name string) bool {
	nodes, err := kube.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("could not get nodes from cluster", "cluster", name, "error", err)
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
