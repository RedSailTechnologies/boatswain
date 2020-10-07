package kraken

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
