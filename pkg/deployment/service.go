package deployment

import (
	"context"

	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/auth"
	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/deployment/template"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
	tw "github.com/redsailtechnologies/boatswain/pkg/twirp"
	pb "github.com/redsailtechnologies/boatswain/rpc/deployment"
	"github.com/redsailtechnologies/boatswain/rpc/repo"
)

var collection = "deployments"

// Service is the implementation for twirp to use
type Service struct {
	auth       auth.Agent
	repo       repo.Repo
	repository *Repository
}

// NewService creates the service
func NewService(a auth.Agent, r repo.Repo, s storage.Storage) *Service {
	return &Service{
		auth:       a,
		repo:       r,
		repository: NewRepository(collection, s),
	}
}

// Create adds a deployment to the list of configurations
func (s Service) Create(ctx context.Context, cmd *pb.CreateDeployment) (*pb.DeploymentCreated, error) {
	if err := s.auth.Authorize(ctx, auth.Admin); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

	d, err := Create(ddd.NewUUID(), cmd.Name, cmd.RepoId, cmd.Branch, cmd.FilePath, ddd.NewTimestamp())
	if err != nil {
		logger.Error("error creating Deployment", "error", err)
		return nil, tw.ToTwirpError(err, "could not create Deployment")
	}

	err = s.repository.Save(d)
	if err != nil {
		logger.Error("error saving Deployment", "error", err)
		return nil, twirp.InternalError("error saving created Deployment")
	}

	return &pb.DeploymentCreated{}, nil
}

// Update edits an already existing deployment
func (s Service) Update(ctx context.Context, cmd *pb.UpdateDeployment) (*pb.DeploymentUpdated, error) {
	if err := s.auth.Authorize(ctx, auth.Admin); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

	d, err := s.repository.Load(cmd.Uuid)
	if err != nil {
		logger.Error("error loading Deployment", "error", err)
		return nil, tw.ToTwirpError(err, "error loading Deployment")
	}

	err = d.Update(cmd.Name, cmd.RepoId, cmd.Branch, cmd.FilePath, ddd.NewTimestamp())
	if err != nil {
		logger.Error("error updating Deployment", "error", err)
		return nil, tw.ToTwirpError(err, "Deployment could not be updated")
	}

	err = s.repository.Save(d)
	if err != nil {
		logger.Error("error saving Deployment", "error", err)
		return nil, twirp.InternalError("error saving updated deployment")
	}

	return &pb.DeploymentUpdated{}, nil
}

// Destroy removes a deployment from the list of configurations
func (s Service) Destroy(ctx context.Context, cmd *pb.DestroyDeployment) (*pb.DeploymentDestroyed, error) {
	if err := s.auth.Authorize(ctx, auth.Admin); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

	d, err := s.repository.Load(cmd.Uuid)
	if err != nil {
		logger.Error("error loading Deployment", "error", err)

		// NOTE we could consider returning the error here
		if err == (ddd.DestroyedError{Entity: "Deployment"}) {
			return &pb.DeploymentDestroyed{}, nil
		}
		return nil, tw.ToTwirpError(err, "error loading Deployment")
	}

	if err = d.Destroy(ddd.NewTimestamp()); err != nil {
		logger.Error("error destroying Deployment", "error", err)
		return nil, tw.ToTwirpError(err, "Deployment could not be destroyed")
	}

	err = s.repository.Save(d)
	if err != nil {
		logger.Error("error saving Deployment", "error", err)
		return nil, twirp.InternalError("error saving destroyed Deployment")
	}

	return &pb.DeploymentDestroyed{}, nil
}

// Read reads out a deployment
func (s Service) Read(ctx context.Context, req *pb.ReadDeployment) (*pb.DeploymentRead, error) {
	if err := s.auth.Authorize(ctx, auth.Reader); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

	d, err := s.repository.Load(req.Uuid)
	if err != nil {
		logger.Error("error reading Deployment", "error", err)
		return nil, tw.ToTwirpError(err, "error loading Deployment")
	}

	r, err := s.repo.Read(ctx, &repo.ReadRepo{Uuid: d.RepoID()})
	if err != nil {
		logger.Error("could not find repo in deployment")
		r = &repo.RepoRead{}
	}

	return &pb.DeploymentRead{
		Uuid:     d.UUID(),
		Name:     d.Name(),
		RepoId:   d.RepoID(),
		RepoName: r.Name,
		Branch:   d.Branch(),
		FilePath: d.FilePath(),
	}, nil
}

// All gets all deployments currently configured and their status
func (s Service) All(ctx context.Context, req *pb.ReadDeployments) (*pb.DeploymentsRead, error) {
	if err := s.auth.Authorize(ctx, auth.Reader); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

	resp := &pb.DeploymentsRead{
		Deployments: make([]*pb.DeploymentReadSummary, 0),
	}

	deployments, err := s.repository.All()
	if err != nil {
		logger.Error("error getting Deployments", "error", err)
		return nil, twirp.InternalError("error loading Deployments")
	}

	for _, d := range deployments {

		resp.Deployments = append(resp.Deployments, &pb.DeploymentReadSummary{
			Uuid: d.UUID(),
			Name: d.Name(),
		})
	}
	return resp, nil
}

// Template takes a deployment and templates its yaml out for verification and viewing
func (s Service) Template(ctx context.Context, req *pb.TemplateDeployment) (*pb.DeploymentTemplated, error) {
	if err := s.auth.Authorize(ctx, auth.Reader); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

	d, err := s.repository.Load(req.Uuid)
	if err != nil {
		logger.Error("error reading Deployment", "error", err)
		return nil, tw.ToTwirpError(err, "error loading Deployment")
	}

	f, err := s.repo.File(ctx, &repo.ReadFile{
		RepoId:   d.RepoID(),
		Branch:   d.Branch(),
		FilePath: d.FilePath(),
	})
	if err != nil {
		logger.Error("could not get deployment file", "error", err)
		return nil, twirp.NotFoundError("deployment file not found")
	}

	e := template.NewEngine(ctx, s.repo)
	yaml, err := e.Template(f.File)
	if err != nil {
		logger.Error("yaml file could not be templated", "error", err)
		return nil, twirp.InternalErrorWith(err)
	}

	return &pb.DeploymentTemplated{
		Yaml: yaml,
	}, nil
}

// Trigger a deployment
func (s Service) Trigger(ctx context.Context, cmd *pb.TriggerDeployment) (*pb.DeploymentTriggered, error) {
	if err := s.auth.Authorize(ctx, auth.Editor); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

	// template/validate the deployment
	d, err := s.repository.Load(cmd.Uuid)
	if err != nil {
		logger.Error("error reading Deployment", "error", err)
		return nil, tw.ToTwirpError(err, "error loading Deployment")
	}

	f, err := s.repo.File(ctx, &repo.ReadFile{
		RepoId:   d.RepoID(),
		Branch:   d.Branch(),
		FilePath: d.FilePath(),
	})
	if err != nil {
		logger.Error("could not get deployment file", "error", err)
		return nil, twirp.NotFoundError("deployment file not found")
	}

	te := template.NewEngine(ctx, s.repo)
	t, err := te.Run(f.File, cmd.Arguments)
	if err != nil {
		logger.Error("yaml file could not be templated", "error", err)
		return nil, twirp.InternalErrorWith(err)
	}
	if err = t.Validate(); err != nil {
		logger.Error("could not validate template", "error", err)
		return nil, twirp.InternalErrorWith(err)
	}

	// validate the trigger
	user := s.auth.User(ctx)
	if err = ValidateTrigger(&user, cmd, *t.Triggers); err != nil {
		logger.Error("invalid trigger", "error", err)
		return nil, twirp.InternalErrorWith(err)
	}

	// create the run and start the engine in the background

	// return the run id

	return nil, nil
}

// Ready implements the ReadyService method so this service can be part of a health check routine
func (s Service) Ready() error {
	return s.repository.store.CheckReady()
}
