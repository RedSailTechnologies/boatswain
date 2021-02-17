package agent

import (
	"context"
	"sync"
	"time"

	"github.com/redsailtechnologies/boatswain/pkg/cluster"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
	pb "github.com/redsailtechnologies/boatswain/rpc/agent"
	"github.com/twitchtv/twirp"
)

// Service is the implementation for twirp to use
type Service struct {
	actions map[string]map[string]*pb.Action
	cluster *cluster.ReadRepository
	results map[string]chan *pb.Result

	actionsLock *sync.Mutex
	resultsLock *sync.Mutex
}

// NewService creates the new service with dependencies
func NewService(s storage.Storage) *Service {
	return &Service{
		actions:     make(map[string]map[string]*pb.Action),
		cluster:     cluster.NewReadRepository(s),
		results:     make(map[string]chan *pb.Result),
		actionsLock: &sync.Mutex{},
		resultsLock: &sync.Mutex{},
	}
}

// Register registers this agent
func (s Service) Register(ctx context.Context, cmd *pb.RegisterAgent) (*pb.AgentRegistered, error) {
	cl, err := s.cluster.Load(cmd.ClusterUuid)
	if err != nil {
		logger.Error("error reading Cluster", "error", err)
		return nil, twirp.InternalError("error loading Cluster")
	}

	return &pb.AgentRegistered{
		ClusterToken: cl.Token(),
	}, nil
}

// Actions gets the next action for the agent or an empty list if there's nothing to do
func (s Service) Actions(ctx context.Context, cmd *pb.ReadActions) (*pb.ActionsRead, error) {
	if cl, err := s.cluster.Load(cmd.ClusterUuid); err != nil || cmd.ClusterToken != cl.Token() {
		return nil, twirp.NewError(twirp.PermissionDenied, "invalid cluster uuid or token")
	}

	if actions, ok := s.actions[cmd.ClusterUuid]; ok {
		actionList := make([]*pb.Action, 0)
		for _, action := range actions {
			actionList = append(actionList, action)
		}
		delete(actions, cmd.ClusterUuid)
		return &pb.ActionsRead{
			Actions: actionList,
		}, nil
	}
	return &pb.ActionsRead{
		Actions: make([]*pb.Action, 0),
	}, nil
}

// Results returns a result for this agent
func (s Service) Results(ctx context.Context, cmd *pb.ReturnResult) (*pb.ResultReturned, error) {
	if cl, err := s.cluster.Load(cmd.ClusterUuid); err != nil || cmd.ClusterToken != cl.Token() {
		return nil, twirp.NewError(twirp.PermissionDenied, "invalid cluster uuid or token")
	}

	s.resultsLock.Lock()
	s.results[cmd.ActionUuid] <- cmd.Result // FIXME
	s.resultsLock.Unlock()

	return &pb.ResultReturned{}, nil
}

// Run runs an action and returns a result from a particular agent
func (s Service) Run(ctx context.Context, cmd *pb.Action) (*pb.Result, error) {
	cluster, err := s.cluster.Load(cmd.ClusterUuid)
	if err != nil || cluster.Token() != cmd.ClusterToken {
		return &pb.Result{
			Data: nil,
		}, err
	}

	s.actionsLock.Lock()
	if _, ok := s.actions[cmd.ClusterUuid]; !ok {
		s.actions[cmd.ClusterUuid] = make(map[string]*pb.Action)
	}
	s.actions[cmd.ClusterUuid][cmd.Uuid] = cmd
	s.actionsLock.Unlock()

	ch := make(chan *pb.Result, 1)
	s.resultsLock.Lock()
	s.results[cmd.Uuid] = ch
	s.resultsLock.Unlock()

	var result *pb.Result
	select {
	case res := <-ch:
		result = res
	case <-time.After(time.Duration(cmd.TimeoutSeconds) * time.Second):
		result = &pb.Result{
			Data:  nil,
			Error: "timeout while waiting for action result",
		}
	}

	s.actionsLock.Lock()
	delete(s.actions[cmd.ClusterUuid], cmd.Uuid)
	s.actionsLock.Unlock()

	s.resultsLock.Lock()
	delete(s.results, cmd.Uuid)
	s.resultsLock.Unlock()

	return result, nil
}
