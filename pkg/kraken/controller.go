package kraken

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClustersController is the controller type for kraken
type ClustersController struct {
	Config *ClusterList
}

// Clusters gets all the clusters currently configured
func (cc *ClustersController) Clusters(c *gin.Context) {
	response := &ClustersResponse{
		Clusters: make([]string, 0),
	}

	for _, cluster := range cc.Config.Clusters {
		response.Clusters = append(response.Clusters, cluster.Name)
	}

	c.JSON(200, response)
}

// Namespaces get all the namespaces for a particular cluster
func (cc *ClustersController) Namespaces(c *gin.Context) {
	var cluster ClusterParam
	if err := c.ShouldBindUri(&cluster); err != nil {
		c.JSON(400, gin.H{"message": "cluster name malformed"})
		return
	}

	config, err := cc.getClusterConfig(cluster.Name)
	if err != nil {
		c.JSON(404, gin.H{"message": "cluster not found"})
		return
	}

	clientset, err := config.ToClientset()
	if err != nil {
		c.JSON(500, gin.H{"message": "error loading cluster configuration"})
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		c.JSON(500, gin.H{"message": "error getting namespaces for cluster"})
		return
	}

	response := make([]string, 0)
	for _, ns := range namespaces.Items {
		response = append(response, ns.Name)
	}

	c.JSON(200, NamespacesResponse{
		Cluster:    cluster.Name,
		Namespaces: response,
	})
	return
}

// Deployments gets all the deployments for a cluster, optionally with a namespace
func (cc *ClustersController) Deployments(c *gin.Context) {
	var cluster ClusterParam
	var namespace NamespaceQuery
	if err := c.ShouldBindUri(&cluster); err != nil {
		c.JSON(400, gin.H{"message": "cluster name malformed"})
		return
	}
	err := c.ShouldBindQuery(&namespace)
	if err != nil {
		c.JSON(400, gin.H{"message": "namespace name malformed"})
		return
	}

	config, err := cc.getClusterConfig(cluster.Name)
	if err != nil {
		c.JSON(404, gin.H{"message": "cluster not found"})
		return
	}

	clientset, err := config.ToClientset()
	if err != nil {
		c.JSON(500, gin.H{"message": "error loading cluster configuration"})
	}

	deployments, err := clientset.AppsV1().Deployments(namespace.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		c.JSON(500, gin.H{"message": "error getting namespaces for cluster"})
		return
	}

	response := &DeploymentsResponse{
		Cluster:     cluster.Name,
		Deployments: make([]DeploymentResponse, 0),
	}
	for _, deploy := range deployments.Items {
		response.Deployments = append(response.Deployments, DeploymentResponse{
			Name:      deploy.ObjectMeta.Name,
			Namespace: deploy.ObjectMeta.Namespace,
			Version:   deploy.ObjectMeta.Labels["version"],
		})
	}

	c.JSON(200, response)
	return
}

func (cc *ClustersController) getClusterConfig(clusterName string) (*Cluster, error) {
	for _, cluster := range cc.Config.Clusters {
		if cluster.Name == clusterName {
			return &cluster, nil
		}
	}
	return &Cluster{}, errors.New("cluster not found")
}
