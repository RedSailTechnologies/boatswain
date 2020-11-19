package repo

import (
	"context"

	"github.com/twitchtv/twirp"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"

	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/helm"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
	pb "github.com/redsailtechnologies/boatswain/rpc/repo"
)

var collection = "repos"

// Service is the implementation for twirp to use
type Service struct {
	helm helm.Agent
	repo *Repository
}

// NewService creates the service
func NewService(h helm.Agent, s storage.Storage) *Service {
	return &Service{
		helm: h,
		repo: NewRepository(collection, s),
	}
}

// Create adds a repo to the list of configurations
func (s Service) Create(ctx context.Context, cmd *pb.CreateRepo) (*pb.RepoCreated, error) {
	r, err := Create(ddd.NewUUID(), cmd.Name, cmd.Endpoint, ddd.NewTimestamp())
	if err != nil {
		logger.Error("error creating Repo", "error", err)
		return nil, toTwirpError(err, "could not create Repo")
	}

	err = s.repo.Save(r)
	if err != nil {
		logger.Error("error saving Repo", "error", err)
		return nil, twirp.InternalError("error saving created repo")
	}

	return &pb.RepoCreated{}, nil
}

// Update edits an already existing repo
func (s Service) Update(ctx context.Context, cmd *pb.UpdateRepo) (*pb.RepoUpdated, error) {
	r, err := s.repo.Load(cmd.Uuid)
	if err != nil {
		logger.Error("error loading cluster", "error", err)
		return nil, toTwirpError(err, "error loading Repo")
	}

	err = r.Update(cmd.Name, cmd.Endpoint, ddd.NewTimestamp())
	if err != nil {
		logger.Error("error updating Repo", "error", err)
		return nil, toTwirpError(err, "Repo could not be updated")
	}

	err = s.repo.Save(r)
	if err != nil {
		logger.Error("error saving Repo", "error", err)
		return nil, twirp.InternalError("error saving updated repo")
	}

	return &pb.RepoUpdated{}, nil
}

// Destroy removes a repo from the list of configurations
func (s Service) Destroy(ctx context.Context, cmd *pb.DestroyRepo) (*pb.RepoDestroyed, error) {
	r, err := s.repo.Load(cmd.Uuid)
	if err != nil {
		logger.Error("error loading Repo", "error", err)
		return nil, toTwirpError(err, "error loading Repo")
	}

	err = r.Destroy(ddd.NewTimestamp())
	if err != nil {
		logger.Error("error destroying Repo", "error", err)
		return nil, toTwirpError(err, "Repo could not be destroyed")
	}

	err = s.repo.Save(r)
	if err != nil {
		logger.Error("error saving Repo", "error", err)
		return nil, twirp.InternalError("error saving destroyed Repo")
	}

	return &pb.RepoDestroyed{}, nil
}

// Read reads out a repo
func (s Service) Read(ctx context.Context, req *pb.ReadRepo) (*pb.RepoRead, error) {
	r, err := s.repo.Load(req.Uuid)
	if err != nil {
		logger.Error("error loading Repo", "error", err)
		return nil, toTwirpError(err, "error loading Repo")
	}

	cr, err := r.toChartRepo()
	if err != nil {
		logger.Error("error converting Repo to helm chart repo", "error", err)
		return nil, twirp.InternalError("error converting Repo to helm chart repo")
	}

	return &pb.RepoRead{
		Uuid:     r.UUID(),
		Name:     r.Name(),
		Endpoint: r.Endpoint(),
		Ready:    s.helm.CheckIndex(cr),
	}, nil
}

// All gets all repos currently configured and their status
func (s Service) All(context.Context, *pb.ReadRepos) (*pb.ReposRead, error) {
	resp := &pb.ReposRead{
		Repos: make([]*pb.RepoRead, 0),
	}

	repos, err := s.repo.All()
	if err != nil {
		logger.Error("error getting Repos", "error", err)
		return nil, twirp.InternalError("error loading Repos")
	}

	for _, r := range repos {
		cr, err := r.toChartRepo()
		if err != nil {
			logger.Error("error converting Repo to helm chart repo", "error", err)
			return nil, twirp.InternalError("error converting Repo to helm chart repo")
		}

		resp.Repos = append(resp.Repos, &pb.RepoRead{
			Uuid:     r.UUID(),
			Name:     r.Name(),
			Endpoint: r.Endpoint(),
			Ready:    s.helm.CheckIndex(cr),
		})
	}
	return resp, nil
}

// Charts gets all the charts for this repository
func (s Service) Charts(ctx context.Context, req *pb.ReadRepo) (*pb.ChartsRead, error) {
	r, err := s.repo.Load(req.Uuid)
	if err != nil {
		logger.Error("error loading Repo", "error", err)
		return nil, toTwirpError(err, "error loading Repo")
	}

	cr, err := r.toChartRepo()
	if err != nil {
		logger.Error("error converting Repo to helm chart repo", "error", err)
		return nil, twirp.InternalError("error converting Repo to helm chart repo")
	}

	chartMap, err := s.helm.GetCharts(cr)
	if err != nil {
		logger.Error("error getting charts", "error", err)
		return nil, twirp.NotFoundError("error getting charts from helm repo")
	}

	charts := make([]*pb.ChartRead, 0)
	for key, val := range chartMap {
		versions := make([]*pb.VersionRead, 0)

		for _, version := range val {
			versions = append(versions, &pb.VersionRead{
				Name:         key,
				ChartVersion: version.Metadata.Version,
				AppVersion:   version.Metadata.AppVersion,
				Description:  version.Metadata.Description,
				Url:          buildChartURL(r.Endpoint(), version.URLs[0]),
			})
		}

		charts = append(charts, &pb.ChartRead{
			Name:     key,
			Versions: versions,
		})
	}

	return &pb.ChartsRead{Charts: charts}, nil
}

func buildChartURL(repoURL string, chart string) string {
	if repoURL[len(repoURL)-1] != byte('/') {
		repoURL = repoURL + "/"
	}
	return repoURL + chart
}

func (r *Repo) toChartRepo() (*repo.ChartRepository, error) {
	providers := []getter.Provider{
		getter.Provider{
			Schemes: []string{"http", "https"},
			New:     getter.NewHTTPGetter,
		},
	}

	entry := &repo.Entry{
		Name: r.Name(),
		URL:  r.Endpoint(),
		// TODO AdamP - we definitely want to support this soon!
		InsecureSkipTLSverify: true,
	}

	return repo.NewChartRepository(entry, providers)
}

func toTwirpError(e error, m string) error {
	switch e.(type) {
	case ddd.DestroyedError:
		return twirp.NotFoundError(e.Error())
	case ddd.InvalidArgumentError:
		return twirp.InvalidArgumentError(e.(ddd.InvalidArgumentError).Arg, e.Error())
	case ddd.NotFoundError:
		return twirp.NotFoundError(e.Error())
	case ddd.RequiredArgumentError:
		return twirp.RequiredArgumentError(e.(ddd.RequiredArgumentError).Arg)
	default:
		return twirp.InternalError(m)
	}
}