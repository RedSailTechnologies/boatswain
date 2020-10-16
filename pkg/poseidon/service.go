package poseidon

import (
	"context"
	"os"

	"github.com/redsailtechnologies/boatswain/pkg/logger"
	pb "github.com/redsailtechnologies/boatswain/rpc/poseidon"
	"github.com/twitchtv/twirp"
)

// Service is the implementation of the Poseidon twirp service
type Service struct {
	config    *Config
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

	return &Service{
		config:    c,
		repoAgent: defaultRepoAgent{},
	}
}

// Charts get all charts and all versions of those charts in the repo
func (s *Service) Charts(ctx context.Context, repo *pb.Repo) (*pb.ChartsResponse, error) {
	cfg, err := s.config.getRepoConfig(repo.Name)
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
	config, err := s.config.getRepoConfig(req.RepoName)
	if err != nil {
		logger.Error("error getting repo config", "error", err)
		return nil, twirp.InternalError("error getting repo config")
	}

	file, err := s.repoAgent.downloadChart(req.ChartName, req.ChartVersion, s.config.CacheDir, config.Endpoint, config.ToChartPathOptions())
	if err != nil {
		logger.Error("error downloading chart", "error", err)
		return nil, twirp.InternalError("error downloading chart")
	}

	return file, nil
}

// Repos gets all helm repos configured
func (s *Service) Repos(ctx context.Context, req *pb.ReposRequest) (*pb.ReposResponse, error) {
	repos := make([]*pb.Repo, 0)
	for _, config := range s.config.Repos {
		helmRepo, err := config.ToChartRepo()
		if err != nil {
			logger.Error("error getting chart repo", "error", err)
			return nil, twirp.InternalError("error getting helm repo")
		}

		repos = append(repos, &pb.Repo{
			Name:     config.Name,
			Endpoint: config.Endpoint,
			Ready:    s.repoAgent.checkIndex(helmRepo),
		})
	}

	return &pb.ReposResponse{Repos: repos}, nil
}

func buildChartURL(repoURL string, chart string) string {
	if repoURL[len(repoURL)-1] != byte('!') {
		repoURL = repoURL + "/"
	}
	return repoURL + chart
}