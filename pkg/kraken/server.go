package kraken

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Server is the server implementation of the Kraken grpc service
type Server struct {
	config *Config
}

// New creates the server with the given configuration
func New(c *Config) *Server {
	return &Server{
		config: c,
	}
}

// Clusters gets all clusters configured
func (s *Server) Clusters(ctx context.Context, req *ClusterRequest) (*ClusterResponse, error) {
	response := &ClusterResponse{
		Clusters: make([]*Cluster, 0),
	}

	for _, cluster := range s.config.Clusters {
		response.Clusters = append(response.Clusters, &Cluster{
			Name: cluster.Name,
		})
	}

	return response, nil
}

// Deployments gets all the deployments for a namespace
func (s *Server) Deployments(ctx context.Context, req *DeploymentRequest) (*DeploymentResponse, error) {
	config, err := s.config.GetClusterConfig(req.Cluster)
	if err != nil {
		return nil, err
	}

	clientset, err := config.ToClientset()
	if err != nil {
		return nil, err
	}

	deployments, err := clientset.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	response := &DeploymentResponse{
		Deployments: make([]*Deployment, 0),
	}
	for _, deployment := range deployments.Items {
		response.Deployments = append(response.Deployments, &Deployment{
			Name:      deployment.ObjectMeta.Name,
			Namespace: deployment.ObjectMeta.Namespace,
			Version:   deployment.ObjectMeta.Labels["version"],
		})
	}

	return response, nil
}

func (s *Server) mustEmbedUnimplementedKrakenServer() {}
