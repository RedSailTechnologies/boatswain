package deployment

import (
	"context"

	"github.com/redsailtechnologies/boatswain/pkg/deployment/template"
	"github.com/redsailtechnologies/boatswain/pkg/deployment/trigger"

	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/auth"
	"github.com/redsailtechnologies/boatswain/pkg/cluster"
	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/deployment/run"
	"github.com/redsailtechnologies/boatswain/pkg/git"
	"github.com/redsailtechnologies/boatswain/pkg/helm"
	"github.com/redsailtechnologies/boatswain/pkg/kube"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/pkg/repo"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
	tw "github.com/redsailtechnologies/boatswain/pkg/twirp"
	pb "github.com/redsailtechnologies/boatswain/rpc/deployment"
)

var collection = "deployments"
var runCollection = "runs"

// Service is the implementation for twirp to use
type Service struct {
	auth    auth.Agent
	git     git.Agent
	cluster *cluster.ReadRepository
	repo    *repo.ReadRepository
	read    *ReadRepository
	write   *writeRepository
	runRead *run.ReadRepository
	store   storage.Storage
	ready   func() error
}

// NewService creates the service
func NewService(a auth.Agent, g git.Agent, s storage.Storage) *Service {
	return &Service{
		auth:    a,
		git:     g,
		cluster: cluster.NewReadRepository(s),
		repo:    repo.NewReadRepository(s),
		read:    NewReadRepository(s),
		write:   newWriteRepository(s),
		runRead: run.NewReadRepository(s),
		store:   s,
		ready:   s.CheckReady,
	}
}

// Create adds a deployment to the list of configurations
func (s Service) Create(ctx context.Context, cmd *pb.CreateDeployment) (*pb.DeploymentCreated, error) {
	// if err := s.auth.Authorize(ctx, auth.Admin); err != nil {
	// 	return nil, tw.ToTwirpError(err, "not authorized")
	// }

	d, err := Create(ddd.NewUUID(), cmd.Name, cmd.RepoId, cmd.Branch, cmd.FilePath, ddd.NewTimestamp())
	if err != nil {
		logger.Error("error creating Deployment", "error", err)
		return nil, tw.ToTwirpError(err, "could not create Deployment")
	}

	err = s.write.save(d)
	if err != nil {
		logger.Error("error saving Deployment", "error", err)
		return nil, twirp.InternalError("error saving created Deployment")
	}

	return &pb.DeploymentCreated{}, nil
}

// Update edits an already existing deployment
func (s Service) Update(ctx context.Context, cmd *pb.UpdateDeployment) (*pb.DeploymentUpdated, error) {
	// if err := s.auth.Authorize(ctx, auth.Admin); err != nil {
	// 	return nil, tw.ToTwirpError(err, "not authorized")
	// }

	d, err := s.read.Load(cmd.Uuid)
	if err != nil {
		logger.Error("error loading Deployment", "error", err)
		return nil, tw.ToTwirpError(err, "error loading Deployment")
	}

	err = d.Update(cmd.Name, cmd.RepoId, cmd.Branch, cmd.FilePath, ddd.NewTimestamp())
	if err != nil {
		logger.Error("error updating Deployment", "error", err)
		return nil, tw.ToTwirpError(err, "Deployment could not be updated")
	}

	err = s.write.save(d)
	if err != nil {
		logger.Error("error saving Deployment", "error", err)
		return nil, twirp.InternalError("error saving updated deployment")
	}

	return &pb.DeploymentUpdated{}, nil
}

// Destroy removes a deployment from the list of configurations
func (s Service) Destroy(ctx context.Context, cmd *pb.DestroyDeployment) (*pb.DeploymentDestroyed, error) {
	// if err := s.auth.Authorize(ctx, auth.Admin); err != nil {
	// 	return nil, tw.ToTwirpError(err, "not authorized")
	// }

	d, err := s.read.Load(cmd.Uuid)
	if err != nil {
		logger.Error("error loading Deployment", "error", err)

		// NOTE we could consider returning the error here
		if err == (ddd.DestroyedError{Entity: entityName}) {
			return &pb.DeploymentDestroyed{}, nil
		}
		return nil, tw.ToTwirpError(err, "error loading Deployment")
	}

	if err = d.Destroy(ddd.NewTimestamp()); err != nil {
		logger.Error("error destroying Deployment", "error", err)
		return nil, tw.ToTwirpError(err, "Deployment could not be destroyed")
	}

	err = s.write.save(d)
	if err != nil {
		logger.Error("error saving Deployment", "error", err)
		return nil, twirp.InternalError("error saving destroyed Deployment")
	}

	return &pb.DeploymentDestroyed{}, nil
}

// Read reads out a deployment
func (s Service) Read(ctx context.Context, req *pb.ReadDeployment) (*pb.DeploymentRead, error) {
	// if err := s.auth.Authorize(ctx, auth.Reader); err != nil {
	// 	return nil, tw.ToTwirpError(err, "not authorized")
	// }

	d, err := s.read.Load(req.Uuid)
	if err != nil {
		logger.Error("error reading Deployment", "error", err)
		return nil, tw.ToTwirpError(err, "error loading Deployment")
	}

	r, err := s.repo.Load(d.RepoID())
	if err != nil {
		logger.Error("could not find repo in deployment", "error", err)
		return nil, twirp.InternalErrorWith(err)
	}

	return &pb.DeploymentRead{
		Uuid:     d.UUID(),
		Name:     d.Name(),
		RepoId:   d.RepoID(),
		RepoName: r.Name(),
		Branch:   d.Branch(),
		FilePath: d.FilePath(),
	}, nil
}

// All gets all deployments currently configured and their status
func (s Service) All(ctx context.Context, req *pb.ReadDeployments) (*pb.DeploymentsRead, error) {
	// if err := s.auth.Authorize(ctx, auth.Reader); err != nil {
	// 	return nil, tw.ToTwirpError(err, "not authorized")
	// }

	resp := &pb.DeploymentsRead{
		Deployments: make([]*pb.DeploymentReadSummary, 0),
	}

	deployments, err := s.read.All()
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
	// if err := s.auth.Authorize(ctx, auth.Reader); err != nil {
	// 	return nil, tw.ToTwirpError(err, "not authorized")
	// }

	d, err := s.read.Load(req.Uuid)
	if err != nil {
		logger.Error("error reading Deployment", "error", err)
		return nil, tw.ToTwirpError(err, "error loading Deployment")
	}

	r, err := s.repo.Load(d.RepoID())
	if err != nil {
		logger.Error("error reading Repo for Deployment", "error", err)
		return nil, tw.ToTwirpError(err, "error loading Repo for Deployment")
	}

	f := s.git.GetFile(r.Endpoint(), r.Token(), d.Branch(), d.FilePath())
	if f == nil {
		logger.Error("could not get deployment file")
		return nil, twirp.NotFoundError("deployment file not found")
	}

	te := template.NewEngine(&git.DefaultAgent{}, s.repo)
	yaml, err := te.Template(f)
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
	// if err := s.auth.Authorize(ctx, auth.Editor); err != nil {
	// 	return nil, tw.ToTwirpError(err, "not authorized")
	// }

	// template/validate the deployment
	d, err := s.read.Load(cmd.Uuid)
	if err != nil {
		logger.Error("error reading Deployment", "error", err)
		return nil, tw.ToTwirpError(err, "error loading Deployment")
	}

	repo, err := s.repo.Load(d.RepoID())
	if err != nil {
		logger.Error("error reading Repo for Deployment", "error", err)
		return nil, tw.ToTwirpError(err, "error loading Repo for Deployment")
	}

	f := s.git.GetFile(repo.Endpoint(), repo.Token(), d.Branch(), d.FilePath())
	if f == nil {
		logger.Error("could not get deployment file")
		return nil, twirp.NotFoundError("deployment file not found")
	}

	te := template.NewEngine(&git.DefaultAgent{}, s.repo)
	t, err := te.Run(f, cmd.Arguments)
	if err != nil {
		logger.Error("yaml file could not be templated", "error", err)
		return nil, twirp.InternalErrorWith(err)
	}
	if err = t.Validate(); err != nil {
		logger.Error("could not validate template", "error", err)
		return nil, twirp.InternalErrorWith(err)
	}

	// validate the trigger
	// user := s.auth.User(ctx)
	// if err = trigger.Validate(&user, cmd, *t); err != nil {
	// 	logger.Error("invalid trigger", "error", err)
	// 	return nil, twirp.InternalErrorWith(err)
	// }

	// create the run
	r, err := run.Create(ddd.NewUUID(), cmd.Uuid, t, &trigger.Trigger{
		Name: cmd.Name,
		Type: getTriggerType(cmd.Type),
		User: auth.User{},
	})

	// start the engine in the background
	eng, err := run.NewEngine(r, s.store, git.DefaultAgent{}, helm.DefaultAgent{}, kube.DefaultAgent{})
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}

	go eng.Run()

	// return the run id
	return &pb.DeploymentTriggered{
		RunUuid: r.UUID(),
	}, nil
}

// Run reads all the information about a particular run
func (s Service) Run(ctx context.Context, req *pb.ReadRun) (*pb.RunRead, error) {
	// if err := s.auth.Authorize(ctx, auth.Reader); err != nil {
	// 	return nil, tw.ToTwirpError(err, "not authorized")
	// }

	r, err := s.runRead.Load(req.DeploymentUuid)
	if err != nil {
		logger.Error("error reading Run", "error", err)
		return nil, tw.ToTwirpError(err, "error loading Run")
	}

	steps := make([]*pb.StepRead, 0)
	for _, step := range r.Steps() {
		steps = append(steps, &pb.StepRead{
			Name:      step.Name,
			Status:    convertStatus(step.Status),
			StartTime: step.Start,
			StopTime:  step.Stop,
			Logs:      convertLogs(step.Logs),
		})
	}

	return &pb.RunRead{
		Uuid:      r.UUID(),
		Version:   r.RunVersion(),
		Status:    convertStatus(r.Status()),
		StartTime: r.StartTime(),
		StopTime:  r.StopTime(),
		Steps:     steps,
	}, nil
}

// Runs reads a summary of all runs for a particular deployment
func (s Service) Runs(ctx context.Context, req *pb.ReadRuns) (*pb.RunsRead, error) {
	// if err := s.auth.Authorize(ctx, auth.Reader); err != nil {
	// 	return nil, tw.ToTwirpError(err, "not authorized")
	// }

	runs, err := s.runRead.All()
	if err != nil {
		logger.Error("error reading Runs", "error", err)
		return nil, tw.ToTwirpError(err, "error loading Run")
	}

	resp := &pb.RunsRead{}
	resp.Runs = make([]*pb.RunReadSummary, 0)

	for _, r := range runs {
		// FIXME - see the RunRepository about optimization here
		if r.DeploymentUUID() == req.DeploymentUuid {
			resp.Runs = append(resp.Runs, &pb.RunReadSummary{
				Uuid:      r.UUID(),
				Version:   r.RunVersion(),
				Status:    convertStatus(r.Status()),
				StartTime: r.StartTime(),
				StopTime:  r.StopTime(),
			})
		}
	}
	return resp, nil

}

// Ready implements the ReadyService method so this service can be part of a health check routine
func (s Service) Ready() error {
	return s.ready()
}

func convertLogs(logs []run.Log) []*pb.StepLog {
	out := make([]*pb.StepLog, 0)
	for _, log := range logs {
		out = append(out, &pb.StepLog{
			Timestamp: log.Timestamp,
			Level:     convertLevel(log.Level),
			Message:   log.Message,
		})
	}
	return out
}

func convertLevel(l run.LogLevel) pb.LogLevel {
	switch l {
	case run.Debug:
		return pb.LogLevel_DEBUG
	case run.Info:
		return pb.LogLevel_INFO
	case run.Warn:
		return pb.LogLevel_WARN
	case run.Error:
		return pb.LogLevel_ERROR
	default:
		return -1
	}
}

func convertStatus(s run.Status) pb.Status {
	switch s {
	case run.NotStarted:
		return pb.Status_NOT_STARTED
	case run.InProgress:
		return pb.Status_IN_PROGRESS
	case run.Failed:
		return pb.Status_FAILED
	case run.Succeeded:
		return pb.Status_SUCCEEDED
	case run.Skipped:
		return pb.Status_SKIPPED
	default:
		return -1
	}
}

func getTriggerType(t pb.TriggerDeployment_TriggerType) trigger.Type {
	switch t {
	case pb.TriggerDeployment_WEB:
		return trigger.WebTrigger
	case pb.TriggerDeployment_MANUAL:
		return trigger.ManualTrigger
	default:
		return trigger.DeploymentTrigger
	}
}
