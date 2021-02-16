package cluster

import (
	"context"

	"github.com/redsailtechnologies/boatswain/pkg/auth"

	"github.com/twitchtv/twirp"

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
	auth  auth.Agent
	k8s   kube.Agent
	read  *ReadRepository
	write *writeRepository
	ready func() error
}

// NewService creates the service
func NewService(a auth.Agent, k kube.Agent, s storage.Storage) *Service {
	return &Service{
		auth:  a,
		k8s:   k,
		read:  NewReadRepository(s),
		write: newWriteRepository(s),
		ready: s.CheckReady,
	}
}

// Create adds a cluster to the list of configurations
func (s Service) Create(ctx context.Context, cmd *pb.CreateCluster) (*pb.ClusterCreated, error) {
	if err := s.auth.Authorize(ctx, auth.Admin); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

	c, err := Create(ddd.NewUUID(), cmd.Name, ddd.NewUUID(), ddd.NewTimestamp())
	if err != nil {
		logger.Error("error creating Cluster", "error", err)
		return nil, tw.ToTwirpError(err, "could not create Cluster")
	}

	err = s.write.save(c)
	if err != nil {
		logger.Error("error saving Cluster", "error", err)
		return nil, twirp.InternalError("error saving created Cluster")
	}

	return &pb.ClusterCreated{}, nil
}

// Update edits an already existing cluster
func (s Service) Update(ctx context.Context, cmd *pb.UpdateCluster) (*pb.ClusterUpdated, error) {
	if err := s.auth.Authorize(ctx, auth.Admin); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

	c, err := s.read.Load(cmd.Uuid)
	if err != nil {
		logger.Error("error loading Cluster", "error", err)
		return nil, tw.ToTwirpError(err, "error loading Cluster")
	}

	err = c.Update(cmd.Name, ddd.NewTimestamp())
	if err != nil {
		logger.Error("error updating Cluster", "error", err)
		return nil, tw.ToTwirpError(err, "Cluster could not be updated")
	}

	err = s.write.save(c)
	if err != nil {
		logger.Error("error saving Cluster", "error", err)
		return nil, twirp.InternalError("error saving updated cluster")
	}

	return &pb.ClusterUpdated{}, nil
}

// Destroy removes a cluster from the list of configurations
func (s Service) Destroy(ctx context.Context, cmd *pb.DestroyCluster) (*pb.ClusterDestroyed, error) {
	if err := s.auth.Authorize(ctx, auth.Admin); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

	c, err := s.read.Load(cmd.Uuid)
	if err != nil {
		logger.Error("error loading Cluster", "error", err)

		// NOTE we could consider returning the error here
		if err == (ddd.DestroyedError{Entity: entityName}) {
			return &pb.ClusterDestroyed{}, nil
		}
		return nil, tw.ToTwirpError(err, "error loading Cluster")
	}

	if err = c.Destroy(ddd.NewTimestamp()); err != nil {
		logger.Error("error destroying Cluster", "error", err)
		return nil, tw.ToTwirpError(err, "Cluster could not be destroyed")
	}

	err = s.write.save(c)
	if err != nil {
		logger.Error("error saving Cluster", "error", err)
		return nil, twirp.InternalError("error saving destroyed Cluster")
	}

	return &pb.ClusterDestroyed{}, nil
}

// Read reads out a cluster
func (s Service) Read(ctx context.Context, req *pb.ReadCluster) (*pb.ClusterRead, error) {
	if err := s.auth.Authorize(ctx, auth.Reader); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

	c, err := s.read.Load(req.Uuid)
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
		Uuid:  c.UUID(),
		Name:  c.Name(),
		Ready: s.k8s.GetClusterStatus(cs, c.Name()), // FIXME
	}, nil
}

// Find finds the cluster uuid by name
func (s Service) Find(ctx context.Context, req *pb.FindCluster) (*pb.ClusterFound, error) {
	if err := s.auth.Authorize(ctx, auth.Reader); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

	clusters, err := s.read.All()
	if err != nil {
		logger.Error("error getting Clusters", "error", err)
		return nil, twirp.InternalError("error loading Clusters")
	}

	for _, cluster := range clusters {
		if cluster.Name() == req.Name {
			return &pb.ClusterFound{
				Uuid: cluster.UUID(),
			}, nil
		}
	}

	return nil, twirp.NotFoundError("could not find repo with that name")
}

// All gets all clusters currently configured and their status
func (s Service) All(ctx context.Context, req *pb.ReadClusters) (*pb.ClustersRead, error) {
	if err := s.auth.Authorize(ctx, auth.Reader); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

	resp := &pb.ClustersRead{
		Clusters: make([]*pb.ClusterRead, 0),
	}

	clusters, err := s.read.All()
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
			Uuid:  c.UUID(),
			Name:  c.Name(),
			Ready: s.k8s.GetClusterStatus(cs, c.Name()),
		})
	}
	return resp, nil
}

// Ready implements the ReadyService method so this service can be part of a health check routine
func (s Service) Ready() error {
	return s.ready()
}
