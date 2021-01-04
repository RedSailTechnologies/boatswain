package health

import (
	"context"

	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/logger"
	pb "github.com/redsailtechnologies/boatswain/rpc/health"
)

// Service is the implementation of a health service
type Service struct {
	services []ReadyService
}

// NewService creates the service with the other services to ready check
func NewService(s ...ReadyService) *Service {
	return &Service{
		services: s,
	}
}

// Live gets if the server is alive
func (s Service) Live(context.Context, *pb.CheckLive) (*pb.LiveCheck, error) {
	// if we've gotten this far, the server is alive
	return &pb.LiveCheck{}, nil
}

// Ready gets if the server is ready for traffic
func (s Service) Ready(ctx context.Context, req *pb.CheckReady) (*pb.ReadyCheck, error) {
	// for ready we wanna do a bit more-the services we get passed in the
	// new method give us a list of services to ready over
	for _, svc := range s.services {
		err := svc.Ready()
		if err != nil {
			logger.Error("ready check failed", "error", err)
			return nil, twirp.InternalError("service not ready")
		}
	}
	return &pb.ReadyCheck{}, nil
}
