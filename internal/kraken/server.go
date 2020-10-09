package kraken

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	pb "github.com/redsailtechnologies/boatswain/pkg/kraken"
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
func (s *Server) Clusters(ctx context.Context, req *pb.ClusterRequest) (*pb.ClusterResponse, error) {
	response := &pb.ClusterResponse{
		Clusters: make([]*pb.Cluster, 0),
	}

	for _, cluster := range s.config.Clusters {
		config, err := s.config.GetClusterConfig(cluster.Name)
		if err != nil {
			return nil, err
		}
		clientset, err := config.ToClientset()
		if err != nil {
			return nil, err
		}

		ready := true
		nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			ready = false
		} else {
			for _, node := range nodes.Items {
				for _, condition := range node.Status.Conditions {
					if condition.Type == "Ready" {
						if condition.Status != "True" {
							ready = false
						}
					}
				}
			}
		}

		response.Clusters = append(response.Clusters, &pb.Cluster{
			Name:     cluster.Name,
			Endpoint: cluster.Endpoint,
			Ready:    ready,
		})
	}

	return response, nil
}

// Deployments gets all the deployments for a namespace
func (s *Server) Deployments(ctx context.Context, req *pb.DeploymentRequest) (*pb.DeploymentResponse, error) {
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

	response := &pb.DeploymentResponse{
		Deployments: make([]*pb.Deployment, 0),
	}
	for _, deployment := range deployments.Items {
		response.Deployments = append(response.Deployments, &pb.Deployment{
			Name:      deployment.ObjectMeta.Name,
			Namespace: deployment.ObjectMeta.Namespace,
			Version:   deployment.ObjectMeta.Labels["version"],
		})
	}

	return response, nil
}

func (s *Server) mustEmbedUnimplementedKrakenServer() {}
