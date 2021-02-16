package agent

import (
	"context"
	"sync"

	"github.com/redsailtechnologies/boatswain/pkg/cluster"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	pb "github.com/redsailtechnologies/boatswain/rpc/agent"
	"github.com/twitchtv/twirp"
)

// Service is the implementation for twirp to use
type Service struct {
	actions map[string][]*pb.Action
	cluster *cluster.ReadRepository
	results map[string]*chan *pb.Result

	actionsLock *sync.Mutex
	resultsLock *sync.Mutex
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
		return &pb.ActionsRead{
			Actions: actions,
		}, nil
	}
	return &pb.ActionsRead{
		Actions: make([]*pb.Action, 0),
	}, nil
}

// Result returns a result for this agent
func (s Service) Result(ctx context.Context, cmd *pb.ReturnResult) (*pb.ResultReturned, error) {
	if cl, err := s.cluster.Load(cmd.ClusterUuid); err != nil || cmd.ClusterToken != cl.Token() {
		return nil, twirp.NewError(twirp.PermissionDenied, "invalid cluster uuid or token")
	}

	*s.results[cmd.ActionUuid] <- cmd.Result

	return &pb.ResultReturned{}, nil
}

// Run runs an action and returns a result from a particular agent
func (s Service) Run(ctx context.Context, cmd *pb.Action) (*pb.Result, error) {
	s.actionsLock.Lock()
	if _, ok := s.actions[cmd.ClusterUuid]; !ok {
		s.actions[cmd.ClusterUuid] = make([]*pb.Action, 0)
	}
	s.actions[cmd.ClusterUuid] = append(s.actions[cmd.ClusterUuid], cmd)
	s.actionsLock.Unlock()

	ch := make(chan *pb.Result, 0)
	s.resultsLock.Lock()
	s.results[cmd.Uuid] = &ch
	s.resultsLock.Unlock()

	result := <-ch
	s.resultsLock.Lock()
	delete(s.results, cmd.Uuid)
	s.resultsLock.Unlock()

	return result, nil
}
