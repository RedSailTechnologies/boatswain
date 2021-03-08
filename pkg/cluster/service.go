package cluster

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/auth"
	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/kube"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
	tw "github.com/redsailtechnologies/boatswain/pkg/twirp"
	"github.com/redsailtechnologies/boatswain/rpc/agent"
	pb "github.com/redsailtechnologies/boatswain/rpc/cluster"
)

var collection = "clusters"

// Service is the implementation for twirp to use
type Service struct {
	agent agent.AgentAction
	auth  auth.Agent
	read  *ReadRepository
	write *writeRepository
	ready func() error
}

// NewService creates the service
func NewService(ag agent.AgentAction, au auth.Agent, s storage.Storage) *Service {
	return &Service{
		agent: ag,
		auth:  au,
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

	return &pb.ClusterCreated{
		Uuid: c.UUID(),
	}, nil
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

	status, err := s.getClusterStatus(c)
	if err != nil {
		logger.Error("error getting Cluster status", "error", err)
	}

	return &pb.ClusterRead{
		Uuid:  c.UUID(),
		Name:  c.Name(),
		Ready: status,
	}, nil
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
		status, err := s.getClusterStatus(c)
		if err != nil {
			logger.Error("error getting Cluster status", "error", err)
		}

		resp.Clusters = append(resp.Clusters, &pb.ClusterRead{
			Uuid:  c.UUID(),
			Name:  c.Name(),
			Ready: status,
		})
	}
	return resp, nil
}

// Token gets the cluster's access token
func (s Service) Token(ctx context.Context, cmd *pb.ReadClusterToken) (*pb.ClusterTokenRead, error) {
	if err := s.auth.Authorize(ctx, auth.Admin); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

	c, err := s.read.Load(cmd.Uuid)
	if err != nil {
		logger.Error("error getting Cluster", "error", err)
		return nil, twirp.NotFoundError("cluster not found")
	}

	return &pb.ClusterTokenRead{
		Token: c.Token(),
	}, nil
}

// Ready implements the ReadyService method so this service can be part of a health check routine
func (s Service) Ready() error {
	return s.ready()
}

func (s Service) getClusterStatus(c *Cluster) (bool, error) {
	args := &kube.Args{}
	jsonArgs, err := json.Marshal(args)
	if err != nil {
		return false, err
	}

	result, err := s.agent.Run(context.Background(), &agent.Action{
		Uuid:           ddd.NewUUID(),
		ClusterUuid:    c.UUID(),
		ClusterToken:   c.Token(),
		ActionType:     agent.ActionType_KUBE_ACTION,
		Action:         string(kube.GetStatus),
		TimeoutSeconds: 2, // FIXME - configurable?
		Args:           jsonArgs,
	})
	if err != nil {
		return false, err
	} else if result.Error != "" {
		return false, errors.New(result.Error)
	}

	return kube.ConvertStatus(result.Data)
}
