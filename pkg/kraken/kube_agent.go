package kraken

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/redsailtechnologies/boatswain/pkg/logger"
)

type kubeAgent interface {
	getClusterDeployments(kubernetes.Interface, string) ([]appsv1.Deployment, error)
	getClusterStatefulSets(kubernetes.Interface, string) ([]appsv1.StatefulSet, error)
	getClusterStatus(kubernetes.Interface, string) bool
}

type defaultKubeAgent struct{}

func (k defaultKubeAgent) getClusterDeployments(kube kubernetes.Interface, name string) ([]appsv1.Deployment, error) {
	d, err := kube.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("could not get deployments from cluster", "cluster", name, "error", err)
		return nil, err
	}
	return d.Items, nil
}

func (k defaultKubeAgent) getClusterStatefulSets(kube kubernetes.Interface, name string) ([]appsv1.StatefulSet, error) {
	ss, err := kube.AppsV1().StatefulSets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("could not get statefulsets from cluster", "cluster", name, "error", err)
		return nil, err
	}
	return ss.Items, nil
}

func (k defaultKubeAgent) getClusterStatus(kube kubernetes.Interface, name string) bool {
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

func toClientset(c *Cluster) (*kubernetes.Clientset, error) {
	restConfig := &rest.Config{
		Host:        c.Endpoint(),
		BearerToken: c.Token(),
		TLSClientConfig: rest.TLSClientConfig{
			CAData: []byte(c.Cert()),
		},
	}
	return kubernetes.NewForConfig(restConfig)
}
