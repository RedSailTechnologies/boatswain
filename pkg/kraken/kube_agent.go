package kraken

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/redsailtechnologies/boatswain/pkg/logger"
	pb "github.com/redsailtechnologies/boatswain/rpc/kraken"
)

// kubeAgent is the agent used to talk to kubernetes clusters
type kubeAgent struct{}

// GetClusterStatus checks a cluster's status by checking each node is ready
func (k *kubeAgent) GetClusterStatus(cs *kubernetes.Clientset, name string) bool {
	nodes, err := cs.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
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

/*
consider the following standard kube labels?
app.kubernetes.io/name: mysql
app.kubernetes.io/instance: mysql-abcxzy
app.kubernetes.io/version: "5.7.21"
app.kubernetes.io/component: database
app.kubernetes.io/part-of: wordpress
app.kubernetes.io/managed-by: helm

or some helm labels?
app.kubernetes.io/name	REC
helm.sh/chart	REC
app.kubernetes.io/managed-by	REC
app.kubernetes.io/instance	REC
app.kubernetes.io/version	OPT
app.kubernetes.io/component	OPT
app.kubernetes.io/part-of	OPT
*/
// deployments
func (k *kubeAgent) GetClusterDeployments(cs *kubernetes.Clientset, cluster *pb.Cluster) ([]*pb.Deployment, error) {
	deployments, err := cs.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	clusterDeployments := make([]*pb.Deployment, 0)
	for _, deployment := range deployments.Items {
		clusterDeployments = append(clusterDeployments, &pb.Deployment{
			Name:      deployment.ObjectMeta.Name,
			Namespace: deployment.ObjectMeta.Namespace,
			Ready:     deploymentIsReady(&deployment.Status),
			Version:   deployment.ObjectMeta.Labels["version"],
			Cluster:   cluster,
		})
	}

	return clusterDeployments, nil
}

func deploymentIsReady(status *v1.DeploymentStatus) bool {
	for _, condition := range status.Conditions {
		if condition.Type == "DeploymentAvailable" {
			return condition.Status == "True"
		}
	}
	return false
}
