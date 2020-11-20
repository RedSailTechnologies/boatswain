package cluster

import (
	"context"

	"github.com/twitchtv/twirp"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/kube"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
	tw "github.com/redsailtechnologies/boatswain/pkg/twirp"
	pb "github.com/redsailtechnologies/boatswain/rpc/cluster"
)

var collection = "clusters"

// Service is the implementation for twirp to use
type Service struct {
	k8s  kube.Agent
	repo *Repository
}

// NewService creates the service
func NewService(k kube.Agent, s storage.Storage) *Service {
	return &Service{
		k8s:  k,
		repo: NewRepository(collection, s),
	}
}

// Create adds a cluster to the list of configurations
func (s Service) Create(ctx context.Context, cmd *pb.CreateCluster) (*pb.ClusterCreated, error) {
	c, err := Create(ddd.NewUUID(), cmd.Name, cmd.Endpoint, cmd.Token, cmd.Cert, ddd.NewTimestamp())
	if err != nil {
		logger.Error("error creating Cluster", "error", err)
		return nil, tw.ToTwirpError(err, "could not create Cluster")
	}

	err = s.repo.Save(c)
	if err != nil {
		logger.Error("error saving Cluster", "error", err)
		return nil, twirp.InternalError("error saving created Cluster")
	}

	return &pb.ClusterCreated{}, nil
}

// Update edits an already existing cluster
func (s Service) Update(ctx context.Context, cmd *pb.UpdateCluster) (*pb.ClusterUpdated, error) {
	c, err := s.repo.Load(cmd.Uuid)
	if err != nil {
		logger.Error("error loading cluster", "error", err)
		return nil, tw.ToTwirpError(err, "error loading Cluster")
	}

	err = c.Update(cmd.Name, cmd.Endpoint, cmd.Token, cmd.Cert, ddd.NewTimestamp())
	if err != nil {
		logger.Error("error updating Cluster", "error", err)
		return nil, tw.ToTwirpError(err, "Cluster could not be updated")
	}

	err = s.repo.Save(c)
	if err != nil {
		logger.Error("error saving Cluster", "error", err)
		return nil, twirp.InternalError("error saving updated cluster")
	}

	return &pb.ClusterUpdated{}, nil
}

// Destroy removes a cluster from the list of configurations
func (s Service) Destroy(ctx context.Context, cmd *pb.DestroyCluster) (*pb.ClusterDestroyed, error) {
	c, err := s.repo.Load(cmd.Uuid)
	if err != nil {
		logger.Error("error loading Cluster", "error", err)

		// NOTE we could consider returning the error here
		if err == (ddd.DestroyedError{Entity: "Cluster"}) {
			return &pb.ClusterDestroyed{}, nil
		}
		return nil, tw.ToTwirpError(err, "error loading Cluster")
	}

	err = c.Destroy(ddd.NewTimestamp())
	if err != nil {
		logger.Error("error destroying Cluster", "error", err)
		return nil, tw.ToTwirpError(err, "Cluster could not be destroyed")
	}

	err = s.repo.Save(c)
	if err != nil {
		logger.Error("error saving Cluster", "error", err)
		return nil, twirp.InternalError("error saving destroyed Cluster")
	}

	return &pb.ClusterDestroyed{}, nil
}

// Read reads out a cluster
func (s Service) Read(ctx context.Context, req *pb.ReadCluster) (*pb.ClusterRead, error) {
	c, err := s.repo.Load(req.Uuid)
	if err != nil {
		logger.Error("error reading Cluster", "error", err)
		return nil, tw.ToTwirpError(err, "error loading Cluster")
	}

	cs, err := c.toClientset()
	if err != nil {
		logger.Error("error converting Cluster to kube clientset", "error", err)
		return nil, twirp.InternalError("error converting Cluster to kubernetes Clientset")
	}

	return &pb.ClusterRead{
		Uuid:     c.UUID(),
		Name:     c.Name(),
		Endpoint: c.Endpoint(),
		Token:    c.Token(),
		Cert:     c.Cert(),
		Ready:    s.k8s.GetClusterStatus(cs, c.Name()),
	}, nil
}

// All gets all clusters currently configured and their status
func (s Service) All(ctx context.Context, req *pb.ReadClusters) (*pb.ClustersRead, error) {
	resp := &pb.ClustersRead{
		Clusters: make([]*pb.ClusterRead, 0),
	}

	clusters, err := s.repo.All()
	if err != nil {
		logger.Error("error getting Clusters", "error", err)
		return nil, twirp.InternalError("error loading Clusters")
	}

	for _, c := range clusters {
		cs, err := c.toClientset()
		if err != nil {
			logger.Error("error converting Cluster to kube clientset", "error", err)
			return nil, twirp.InternalError("error converting Cluster to kubernetes clientset")
		}

		resp.Clusters = append(resp.Clusters, &pb.ClusterRead{
			Uuid:     c.UUID(),
			Name:     c.Name(),
			Endpoint: c.Endpoint(),
			Token:    c.Token(),
			Cert:     c.Cert(),
			Ready:    s.k8s.GetClusterStatus(cs, c.Name()),
		})
	}
	return resp, nil
}

func (c *Cluster) toClientset() (*kubernetes.Clientset, error) {
	restConfig := &rest.Config{
		Host:        c.Endpoint(),
		BearerToken: c.Token(),
		TLSClientConfig: rest.TLSClientConfig{
			CAData: []byte(c.Cert()),
		},
	}
	return kubernetes.NewForConfig(restConfig)
}
