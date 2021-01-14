package deployment

import (
	"context"

	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/auth"
	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/git"
	"github.com/redsailtechnologies/boatswain/pkg/helm"
	"github.com/redsailtechnologies/boatswain/pkg/kube"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
	tw "github.com/redsailtechnologies/boatswain/pkg/twirp"
	"github.com/redsailtechnologies/boatswain/rpc/cluster"
	pb "github.com/redsailtechnologies/boatswain/rpc/deployment"
	"github.com/redsailtechnologies/boatswain/rpc/repo"
)

var collection = "deployments"
var runCollection = "runs"

// Service is the implementation for twirp to use
type Service struct {
	auth          auth.Agent
	cluster       cluster.Cluster
	repo          repo.Repo
	repository    *Repository
	runRepository *RunRepository
}

// NewService creates the service
func NewService(a auth.Agent, c cluster.Cluster, r repo.Repo, s storage.Storage) *Service {
	return &Service{
		auth:          a,
		cluster:       c,
		repo:          r,
		repository:    NewRepository(collection, s),
		runRepository: NewRunRepository(runCollection, s),
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

	te := NewTemplateEngine(ctx, s.repo)
	yaml, err := te.Template(f.File)
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

	// get all clusters and repos for use later
	repos, err := s.repo.All(ctx, &repo.ReadRepos{})
	if err != nil {
		logger.Error("could not get repos", "error", err)
		return nil, err
	}
	clusters, err := s.cluster.All(ctx, &cluster.ReadClusters{})
	if err != nil {
		logger.Error("could not get clusters", "error", err)
		return nil, err
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

	te := NewTemplateEngine(ctx, s.repo)
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
	if err = ValidateTrigger(&user, cmd, *t); err != nil {
		logger.Error("invalid trigger", "error", err)
		return nil, twirp.InternalErrorWith(err)
	}

	// create the run
	run, err := CreateRun(ddd.NewUUID(), cmd.Uuid, ddd.NewTimestamp(), t, &Trigger{
		Name: cmd.Name,
		Type: getTriggerType(cmd.Type),
		User: user,
	})
	if err = s.runRepository.Save(run); err != nil {
		logger.Error("could not save created run", "error", err)
		return nil, twirp.InternalErrorWith(err)
	}

	// start the engine in the background
	agents := runAgents{
		git:  git.DefaultAgent{},
		helm: helm.DefaultAgent{},
		kube: kube.DefaultAgent{},
	}
	entities := runEntities{
		clusters.Clusters,
		repos.Repos,
	}
	re := NewRunEngine(run, s.runRepository, agents, entities)

	go re.Run()

	// return the run id
	return &pb.DeploymentTriggered{
		RunUuid: run.UUID(),
	}, nil
}

// Run reads all the information about a particular run
func (s Service) Run(ctx context.Context, req *pb.ReadRun) (*pb.RunRead, error) {
	if err := s.auth.Authorize(ctx, auth.Reader); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

	r, err := s.runRepository.Load(req.DeploymentUuid)
	if err != nil {
		logger.Error("error reading Run", "error", err)
		return nil, tw.ToTwirpError(err, "error loading Run")
	}

	steps := make([]*pb.StepRead, 0)
	for _, step := range *r.Strategy {
		steps = append(steps, &pb.StepRead{
			Name:   step.Name,
			Status: convertStatus(step.Status),
			Log:    string(step.Log),
		})
	}

	return &pb.RunRead{
		Uuid:    r.UUID(),
		Version: r.RunVersion(),
		Status:  convertStatus(r.Status()),
		Steps:   steps,
	}, nil
}

// Runs reads a summary of all runs for a particular deployment
func (s Service) Runs(ctx context.Context, req *pb.ReadRuns) (*pb.RunsRead, error) {
	if err := s.auth.Authorize(ctx, auth.Reader); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

	runs, err := s.runRepository.All()
	if err != nil {
		logger.Error("error reading Run", "error", err)
		return nil, tw.ToTwirpError(err, "error loading Run")
	}

	resp := &pb.RunsRead{}
	resp.Runs = make([]*pb.RunReadSummary, 0)

	for _, r := range runs {
		// FIXME - see the RunRepository about optimization here
		if r.DeploymentUUID() == req.DeploymentUuid {
			resp.Runs = append(resp.Runs, &pb.RunReadSummary{
				Uuid:    r.UUID(),
				Version: r.RunVersion(),
				Status:  convertStatus(r.Status()),
			})
		}
	}
	return resp, nil

}

// Ready implements the ReadyService method so this service can be part of a health check routine
func (s Service) Ready() error {
	return s.repository.store.CheckReady()
}

func convertStatus(s Status) pb.Status {
	switch s {
	case NotStarted:
		return pb.Status_NotStarted
	case InProgress:
		return pb.Status_InProgress
	case Failed:
		return pb.Status_Failed
	case Succeeded:
		return pb.Status_Succeeded
	case Skipped:
		return pb.Status_Skipped
	default:
		return -1
	}
}

func getTriggerType(t pb.TriggerDeployment_TriggerType) TriggerType {
	switch t {
	case pb.TriggerDeployment_WEB:
		return WebTrigger
	case pb.TriggerDeployment_MANUAL:
		return ManualTrigger
	default:
		return DeploymentTrigger
	}
}
