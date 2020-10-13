package kraken

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/redsailtechnologies/boatswain/pkg/logger"
)

type kubeAgent interface {
	getClusterStatus(kube kubernetes.Interface, name string) bool
}

// kubeAgent is the agent used to talk to kubernetes clusters
type defaultKubeAgent struct{}

// GetClusterStatus checks a cluster's status by checking each node is ready
func (k defaultKubeAgent) getClusterStatus(kube kubernetes.Interface, name string) bool {
	nodes, err := kube.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Warn("could not get nodes from cluster", "cluster", name)
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
