package poseidon

import (
	"context"
	"errors"
	"os"

	"github.com/google/uuid"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	pb "github.com/redsailtechnologies/boatswain/rpc/poseidon"
	"github.com/twitchtv/twirp"
)

// Service is the implementation of the Poseidon twirp service
type Service struct {
	cacheDir  string
	repos     []*Repo
	repoAgent repoAgent
}

// New creates the Service with the given configuraiton
func New(c *Config) *Service {
	if _, err := os.Stat(c.CacheDir); os.IsNotExist(err) {
		logger.Info("creating directory for chart cache", "dir", c.CacheDir)
		if err := os.Mkdir(c.CacheDir, 0700); err != nil {
			logger.Fatal("could not create chart cache directory", "error", err)
		}
	} else if err != nil {
		logger.Fatal("could not check existence of chart cache directory", "error", err)
	}

	repos := make([]*Repo, len(c.Repos))
	for i, repo := range c.Repos {
		repos[i] = &Repo{
			&pb.Repo{
				Uuid:     uuid.New().String(),
				Name:     repo.Name,
				Endpoint: repo.Endpoint,
			},
		}
	}

	return &Service{
		cacheDir:  c.CacheDir,
		repos:     repos,
		repoAgent: defaultRepoAgent{},
	}
}

// Charts get all charts and all versions of those charts in the repo
func (s *Service) Charts(ctx context.Context, repo *pb.Repo) (*pb.ChartsResponse, error) {
	cfg, err := s.getRepoConfig(repo.Name)
	if err != nil {
		logger.Error("error getting repo config", "error", err)
		return nil, twirp.InternalError("error getting repo config")
	}

	helmRepo, err := cfg.ToChartRepo()
	if err != nil {
		logger.Error("error getting chart repo", "error", err)
		return nil, twirp.InternalError("error getting helm repo")
	}

	chartMap, err := s.repoAgent.getCharts(helmRepo)
	if err != nil {
		logger.Error("error getting charts", "error", err)
		return nil, twirp.NotFoundError("error getting charts from helm repo")
	}

	charts := make([]*pb.Chart, 0)
	for key, val := range chartMap {
		versions := make([]*pb.ChartVersion, 0)

		for _, version := range val {
			versions = append(versions, &pb.ChartVersion{
				Name:         key,
				ChartVersion: version.Metadata.Version,
				AppVersion:   version.Metadata.AppVersion,
				Description:  version.Metadata.Description,
				Url:          buildChartURL(repo.Endpoint, version.URLs[0]),
			})
		}

		charts = append(charts, &pb.Chart{
			Name:     key,
			Versions: versions,
		})
	}

	return &pb.ChartsResponse{Charts: charts}, nil
}

// DownloadChart gets the chart file from the repo
func (s *Service) DownloadChart(ctx context.Context, req *pb.DownloadRequest) (*pb.File, error) {
	config, err := s.getRepoConfig(req.RepoName)
	if err != nil {
		logger.Error("error getting repo config", "error", err)
		return nil, twirp.InternalError("error getting repo config")
	}

	file, err := s.repoAgent.downloadChart(req.ChartName, req.ChartVersion, s.cacheDir, config.Endpoint, config.ToChartPathOptions())
	if err != nil {
		logger.Error("error downloading chart", "error", err)
		return nil, twirp.InternalError("error downloading chart")
	}

	return file, nil
}

// AddRepo adds a repo to this service
func (s *Service) AddRepo(ctx context.Context, repo *pb.Repo) (*pb.EmptyResponse, error) {
	if !testRepoURL(repo.Endpoint) {
		return nil, twirp.InvalidArgumentError("Repo.Endpoint", "repo url must begin with http:// or https://")
	}

	s.repos = append(s.repos, &Repo{
		&pb.Repo{
			Uuid:     uuid.New().String(),
			Name:     repo.Name,
			Endpoint: repo.Endpoint,
		},
	})
	return &pb.EmptyResponse{}, nil
}

// DeleteRepo deletes a repo from this service
func (s *Service) DeleteRepo(ctx context.Context, repo *pb.Repo) (*pb.EmptyResponse, error) {
	for i := range s.repos {
		if s.repos[i].Uuid == repo.Uuid {
			s.repos[i] = s.repos[len(s.repos)-1]
			s.repos = s.repos[:len(s.repos)-1]
			return &pb.EmptyResponse{}, nil
		}
	}
	return nil, twirp.InternalError("repo not found")
}

// EditRepo edits a repo in this service
func (s *Service) EditRepo(ctx context.Context, repo *pb.Repo) (*pb.EmptyResponse, error) {
	if !testRepoURL(repo.Endpoint) {
		return nil, twirp.InvalidArgumentError("Repo.Endpoint", "repo url must begin with http:// or https://")
	}

	for i := range s.repos {
		if s.repos[i].Uuid == repo.Uuid {
			s.repos[i] = &Repo{repo}
			return &pb.EmptyResponse{}, nil
		}
	}
	return nil, twirp.InternalError("repo not found")
}

// Repos gets all helm repos configured
func (s *Service) Repos(ctx context.Context, req *pb.ReposRequest) (*pb.ReposResponse, error) {
	repos := make([]*pb.Repo, 0)
	for _, repo := range s.repos {
		helmRepo, err := repo.ToChartRepo()
		if err != nil {
			logger.Error("error getting chart repo", "error", err)
			return nil, twirp.InternalError("error getting helm repo")
		}

		repos = append(repos, &pb.Repo{
			Uuid:     repo.Uuid,
			Name:     repo.Name,
			Endpoint: repo.Endpoint,
			Ready:    s.repoAgent.checkIndex(helmRepo),
		})
	}

	return &pb.ReposResponse{Repos: repos}, nil
}

func (s *Service) getRepoConfig(repoName string) (*Repo, error) {
	for _, repo := range s.repos {
		if repo.Name == repoName {
			return repo, nil
		}
	}
	return nil, errors.New("repo not found")
}

func buildChartURL(repoURL string, chart string) string {
	if repoURL[len(repoURL)-1] != byte('/') {
		repoURL = repoURL + "/"
	}
	return repoURL + chart
}

func testRepoURL(repoURL string) bool {
	if len(repoURL) < 7 {
		return false
	} else if repoURL[:7] == "http://" {
		return true
	}
	if len(repoURL) < 8 || repoURL[:8] != "https://" {
		return false
	}
	return true
}
