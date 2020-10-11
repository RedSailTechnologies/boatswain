package kraken

import (
	"context"
	"errors"

	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/logger"
	pb "github.com/redsailtechnologies/boatswain/rpc/kraken"
)

// Service is the implementation of the Kraken grpc service
type Service struct {
	config *Config
	agent  *kubeAgent
}

// New creates the Service with the given configuration
func New(c *Config) *Service {
	return &Service{
		config: c,
		agent:  &kubeAgent{},
	}
}

// Clusters gets all clusters configured
func (s *Service) Clusters(ctx context.Context, req *pb.ClustersRequest) (*pb.ClustersResponse, error) {
	response := &pb.ClustersResponse{
		Clusters: make([]*pb.Cluster, 0),
	}

	for _, cluster := range s.config.Clusters {
		clientset, err := s.config.ToClientset(cluster.Name)
		if err != nil {
			logger.Error("could not get clientset for cluster", "cluster", cluster.Name)
			return nil, twirp.InternalError("error getting cluster clientset")
		}

		response.Clusters = append(response.Clusters, &pb.Cluster{
			Name:     cluster.Name,
			Endpoint: cluster.Endpoint,
			Ready:    s.agent.GetClusterStatus(clientset, cluster.Name),
		})
	}

	return response, nil
}

// ClusterStatus gets the status of a cluster
func (s *Service) ClusterStatus(ctx context.Context, cluster *pb.Cluster) (*pb.Cluster, error) {
	return nil, errors.New("not implemented")
}

// Deployments gets all the deployments for a namespace
func (s *Service) Deployments(ctx context.Context, req *pb.DeploymentsRequest) (*pb.DeploymentsResponse, error) {
	clientset, err := s.config.ToClientset(req.Cluster.Name)
	if err != nil {
		return nil, err
	}

	deployments, err := s.agent.GetClusterDeployments(clientset, req.Cluster)
	if err != nil {
		return nil, err
	}

	return &pb.DeploymentsResponse{
		Deployments: deployments,
	}, nil
}

// DeploymentStatus gets the status for a single deployment
func (s *Service) DeploymentStatus(ctx context.Context, deployment *pb.Deployment) (*pb.Deployment, error) {
	return nil, errors.New("not implemented")
}
