package kraken

import (
	"github.com/gin-gonic/gin"
)

// ClustersController is the controller type for kraken
type ClustersController struct{}

// Clusters gets all the clusters currently configured
func (cc *ClustersController) Clusters(c *gin.Context) {
	c.JSON(200, ClustersResponse{
		Clusters: []string{"cluster-a", "cluster-b"},
	})
}

// Namespaces get all the namespaces for a particular cluster
func (cc *ClustersController) Namespaces(c *gin.Context) {
	var cluster ClusterParam
	if err := c.ShouldBindUri(&cluster); err != nil {
		c.JSON(400, gin.H{"message": "cluster name incorrect"})
	}

	c.JSON(200, NamespacesResponse{
		Cluster:    cluster.Name,
		Namespaces: []string{"namespace-a", "namespace-b", "namespace-c"},
	})
}

// Deployments gets all the deployments for a cluster, optionally with a namespace
func (cc *ClustersController) Deployments(c *gin.Context) {
	var cluster ClusterParam
	var namespace NamespaceQuery
	if err := c.ShouldBindUri(&cluster); err != nil {
		c.JSON(400, gin.H{"message": "cluster name incorrect"})
	}
	err := c.ShouldBindQuery(&namespace)
	if err == nil && namespace.Namespace != "" {
		c.JSON(200, DeploymentsResponse{
			Cluster: cluster.Name,
			Deployments: []DeploymentResponse{
				DeploymentResponse{
					Name:      "deployment-a",
					Namespace: namespace.Namespace,
					Version:   "0.1.0",
				},
			},
		})
	} else {
		c.JSON(200, DeploymentsResponse{
			Cluster: cluster.Name,
			Deployments: []DeploymentResponse{
				DeploymentResponse{
					Name:      "deployment-a",
					Namespace: "namespace-a",
					Version:   "0.1.0",
				},
			},
		})
	}
}

// ClustersResponse is the response given to a call to get the configured clusters
type ClustersResponse struct {
	Clusters []string `json:"clusters" binding:"required"`
}

// ClusterParam is the uri query for the cluster
type ClusterParam struct {
	Name string `uri:"name" binding:"required"`
}

// NamespaceQuery is the query param for the namespace
type NamespaceQuery struct {
	Namespace string `form:"namespace"`
}

// NamespacesResponse is the resposne given to a call to get the namespaces in a cluster
type NamespacesResponse struct {
	Cluster    string   `json:"cluster" binding:"required"`
	Namespaces []string `json:"namespaces" binding:"required"`
}

// DeploymentsResponse is the response given to a call to get the deployments
type DeploymentsResponse struct {
	Cluster     string               `json:"cluster" binding:"required"`
	Deployments []DeploymentResponse `json:"deployments" binding:"required"`
}

// DeploymentResponse is a single deployment from a deployments query
type DeploymentResponse struct {
	Name      string `json:"name" binding:"required"`
	Namespace string `json:"namespace" binding:"required"`
	Version   string `json:"version"`
}
