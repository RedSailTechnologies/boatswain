package poseidon

import (
	"context"
	"errors"

	"github.com/redsailtechnologies/boatswain/pkg/logger"
	pb "github.com/redsailtechnologies/boatswain/rpc/poseidon"
)

// Service is the implementation of the Poseidon twirp service
type Service struct {
	config    *Config
	repoAgent repoAgent
}

// New creates the Service with the given configuraiton
func New(c *Config) *Service {
	return &Service{
		config:    c,
		repoAgent: defaultRepoAgent{},
	}
}

// Repos gets all helm repos configured
func (s *Service) Repos(ctx context.Context, req *pb.ReposRequest) (*pb.ReposResponse, error) {
	repos := make([]*pb.Repo, 0)
	for _, config := range s.config.Repos {
		helmRepo, err := config.ToChartRepo()
		if err != nil {
			logger.Error("error getting chart repo", "error", err)
			return nil, errors.New("error getting chart repo")
		}

		repos = append(repos, &pb.Repo{
			Name:     config.Name,
			Endpoint: config.Endpoint,
			Ready:    s.repoAgent.checkIndex(helmRepo),
		})
	}

	return &pb.ReposResponse{Repos: repos}, nil
}
