package repo

import (
	"context"

	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/auth"
	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/git"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
	tw "github.com/redsailtechnologies/boatswain/pkg/twirp"
	pb "github.com/redsailtechnologies/boatswain/rpc/repo"
)

var collection = "repos"

// Service is the implementation for twirp to use
type Service struct {
	auth  auth.Agent
	git   git.Agent
	repo  Agent
	read  *ReadRepository
	write *writeRepository
	ready func() error
}

// NewService creates the service
func NewService(a auth.Agent, g git.Agent, r Agent, s storage.Storage) *Service {
	return &Service{
		auth:  a,
		git:   g,
		repo:  r,
		read:  NewReadRepository(s),
		write: newWriteRepository(s),
		ready: s.CheckReady,
	}
}

// Create adds a repo to the list of configurations
func (s Service) Create(ctx context.Context, cmd *pb.CreateRepo) (*pb.RepoCreated, error) {
	if err := s.auth.Authorize(ctx, auth.Admin); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

	r, err := Create(ddd.NewUUID(), cmd.Name, cmd.Endpoint, cmd.Token, Type(cmd.Type), ddd.NewTimestamp())
	if err != nil {
		logger.Error("error creating Repo", "error", err)
		return nil, tw.ToTwirpError(err, "could not create Repo")
	}

	err = s.write.save(r)
	if err != nil {
		logger.Error("error saving Repo", "error", err)
		return nil, twirp.InternalError("error saving created repo")
	}

	return &pb.RepoCreated{}, nil
}

// Update edits an already existing repo
func (s Service) Update(ctx context.Context, cmd *pb.UpdateRepo) (*pb.RepoUpdated, error) {
	if err := s.auth.Authorize(ctx, auth.Admin); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

	r, err := s.read.Load(cmd.Uuid)
	if err != nil {
		logger.Error("error loading repo", "error", err)
		return nil, tw.ToTwirpError(err, "error loading Repo")
	}

	err = r.Update(cmd.Name, cmd.Endpoint, cmd.Token, Type(cmd.Type), ddd.NewTimestamp())
	if err != nil {
		logger.Error("error updating Repo", "error", err)
		return nil, tw.ToTwirpError(err, "Repo could not be updated")
	}

	err = s.write.save(r)
	if err != nil {
		logger.Error("error saving Repo", "error", err)
		return nil, twirp.InternalError("error saving updated repo")
	}

	return &pb.RepoUpdated{}, nil
}

// Destroy removes a repo from the list of configurations
func (s Service) Destroy(ctx context.Context, cmd *pb.DestroyRepo) (*pb.RepoDestroyed, error) {
	if err := s.auth.Authorize(ctx, auth.Admin); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

	r, err := s.read.Load(cmd.Uuid)
	if err != nil {
		logger.Error("error loading Repo", "error", err)

		// NOTE we could consider returning the error here
		if err == (ddd.DestroyedError{Entity: entityName}) {
			return &pb.RepoDestroyed{}, nil
		}
		return nil, tw.ToTwirpError(err, "error loading Repo")
	}

	if err = r.Destroy(ddd.NewTimestamp()); err != nil {
		logger.Error("error destroying Repo", "error", err)
		return nil, tw.ToTwirpError(err, "Repo could not be destroyed")
	}

	err = s.write.save(r)
	if err != nil {
		logger.Error("error saving Repo", "error", err)
		return nil, twirp.InternalError("error saving destroyed Repo")
	}

	return &pb.RepoDestroyed{}, nil
}

// Read reads out a repo
func (s Service) Read(ctx context.Context, req *pb.ReadRepo) (*pb.RepoRead, error) {
	if err := s.auth.Authorize(ctx, auth.Reader); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

	r, err := s.read.Load(req.Uuid)
	if err != nil {
		logger.Error("error loading Repo", "error", err)
		return nil, tw.ToTwirpError(err, "error loading Repo")
	}

	var ready bool
	switch r.Type() {
	case HELM:
		ready = s.helmRepoReady(r)
	case GIT:
		ready = s.gitRepoReady(r)
	default:
		ready = false
	}

	return &pb.RepoRead{
		Uuid:     r.UUID(),
		Name:     r.Name(),
		Endpoint: r.Endpoint(),
		Type:     pb.RepoType(r.Type()),
		Ready:    ready,
	}, nil
}

// Find finds the repo uuid by name
func (s Service) Find(ctx context.Context, req *pb.FindRepo) (*pb.RepoFound, error) {
	if err := s.auth.Authorize(ctx, auth.Reader); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

	repos, err := s.read.All()
	if err != nil {
		logger.Error("error getting Repos", "error", err)
		return nil, twirp.InternalError("error loading Repos")
	}

	for _, repo := range repos {
		if repo.Name() == req.Name {
			return &pb.RepoFound{
				Uuid: repo.UUID(),
			}, nil
		}
	}

	return nil, twirp.NotFoundError("could not find repo with that name")
}

// All gets all repos currently configured and their status
func (s Service) All(ctx context.Context, req *pb.ReadRepos) (*pb.ReposRead, error) {
	if err := s.auth.Authorize(ctx, auth.Reader); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

	resp := &pb.ReposRead{
		Repos: make([]*pb.RepoRead, 0),
	}

	repos, err := s.read.All()
	if err != nil {
		logger.Error("error getting Repos", "error", err)
		return nil, twirp.InternalError("error loading Repos")
	}

	for _, r := range repos {
		var ready bool
		switch r.Type() {
		case HELM:
			ready = s.helmRepoReady(r)
		case GIT:
			ready = s.gitRepoReady(r)
		default:
			ready = false
		}

		resp.Repos = append(resp.Repos, &pb.RepoRead{
			Uuid:     r.UUID(),
			Name:     r.Name(),
			Endpoint: r.Endpoint(),
			Type:     pb.RepoType(r.Type()),
			Ready:    ready,
		})
	}
	return resp, nil
}

// File gets the contents of a file from the git repo
func (s Service) File(ctx context.Context, req *pb.ReadFile) (*pb.FileRead, error) {
	if err := s.auth.Authorize(ctx, auth.Editor); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

	r, err := s.read.Load(req.RepoId)
	if err != nil {
		return nil, tw.ToTwirpError(err, "error loading Repo")
	}

	if r.Type() != GIT {
		logger.Warn("cannot get file for non-git repo")
		return nil, twirp.InvalidArgumentError("repo", "must be a git repo")
	}

	if !s.gitRepoReady(r) {
		logger.Warn("cannot get file from git repo with offline status")
		return nil, twirp.InvalidArgumentError("repo", "status is offline")
	}

	file := s.git.GetFile(r.Endpoint(), r.Token(), req.Branch, req.FilePath)
	if file == nil {
		return nil, twirp.NotFoundError("file not found")
	}

	return &pb.FileRead{
		File: file,
	}, nil
}

// Ready implements the ReadyService method so this service can be part of a health check routine
func (s Service) Ready() error {
	return s.ready()
}

func (s Service) gitRepoReady(r *Repo) bool {
	return s.git.CheckRepo(r.Endpoint(), r.Token())
}

func (s Service) helmRepoReady(r *Repo) bool {
	return s.repo.CheckIndex(r.Name(), r.Endpoint(), r.Token())
}

func buildChartURL(repoURL string, chart string) string {
	if repoURL[len(repoURL)-1] != byte('/') {
		repoURL = repoURL + "/"
	}
	return repoURL + chart
}
