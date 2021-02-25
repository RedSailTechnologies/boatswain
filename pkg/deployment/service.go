package deployment

import (
	"context"
	"errors"

	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/auth"
	"github.com/redsailtechnologies/boatswain/pkg/cluster"
	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/deployment/approval"
	"github.com/redsailtechnologies/boatswain/pkg/deployment/run"
	"github.com/redsailtechnologies/boatswain/pkg/deployment/template"
	"github.com/redsailtechnologies/boatswain/pkg/deployment/trigger"
	"github.com/redsailtechnologies/boatswain/pkg/git"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/pkg/repo"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
	tw "github.com/redsailtechnologies/boatswain/pkg/twirp"
	"github.com/redsailtechnologies/boatswain/rpc/agent"
	pb "github.com/redsailtechnologies/boatswain/rpc/deployment"
	tr "github.com/redsailtechnologies/boatswain/rpc/trigger"
)

var collection = "deployments"
var runCollection = "runs"

// Service is the implementation for twirp to use
type Service struct {
	auth      auth.Agent
	agent     agent.AgentAction
	git       git.Agent
	cluster   *cluster.ReadRepository
	repo      *repo.ReadRepository
	read      *ReadRepository
	write     *writeRepository
	runRead   *run.ReadRepository
	aprvRead  *approval.ReadRepository
	aprvWrite *approval.WriteRepository
	store     storage.Storage
	ready     func() error
}

// NewService creates the service
func NewService(ag agent.AgentAction, au auth.Agent, g git.Agent, s storage.Storage) *Service {
	return &Service{
		agent:     ag,
		auth:      au,
		git:       g,
		cluster:   cluster.NewReadRepository(s),
		repo:      repo.NewReadRepository(s),
		read:      NewReadRepository(s),
		write:     newWriteRepository(s),
		runRead:   run.NewReadRepository(s),
		aprvRead:  approval.NewReadRepository(s),
		aprvWrite: approval.NewWriteRepository(s),
		store:     s,
		ready:     s.CheckReady,
	}
}

// Create adds a deployment to the list of configurations
func (s Service) Create(ctx context.Context, cmd *pb.CreateDeployment) (*pb.DeploymentCreated, error) {
	if err := s.auth.Authorize(ctx, auth.Admin); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

	d, err := Create(ddd.NewUUID(), cmd.Name, ddd.NewUUID(), cmd.RepoId, cmd.Branch, cmd.FilePath, ddd.NewTimestamp())
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
	if err := s.auth.Authorize(ctx, auth.Admin); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

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
	if err := s.auth.Authorize(ctx, auth.Admin); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

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
	if err := s.auth.Authorize(ctx, auth.Reader); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

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
	if err := s.auth.Authorize(ctx, auth.Reader); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

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
	if err := s.auth.Authorize(ctx, auth.Reader); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

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

// Token gets the token for this deployment, for use with web calls
func (s Service) Token(ctx context.Context, req *pb.ReadToken) (*pb.TokenRead, error) {
	if err := s.auth.Authorize(ctx, auth.Admin); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

	d, err := s.read.Load(req.Uuid)
	if err != nil {
		logger.Error("error reading Deployment", "error", err)
		return nil, tw.ToTwirpError(err, "error loading Deployment")
	}

	return &pb.TokenRead{
		Token: d.Token(),
	}, nil
}

// Run reads all the information about a particular run
func (s Service) Run(ctx context.Context, req *pb.ReadRun) (*pb.RunRead, error) {
	if err := s.auth.Authorize(ctx, auth.Reader); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

	r, err := s.runRead.Load(req.DeploymentUuid)
	if err != nil {
		logger.Error("error reading Run", "error", err)
		return nil, tw.ToTwirpError(err, "error loading Run")
	}

	steps := make([]*pb.StepRead, len(r.Steps()))
	for i, step := range r.Steps() {
		steps[i] = &pb.StepRead{
			Name:      step.Name,
			Status:    convertStatus(step.Status),
			StartTime: step.Start,
			StopTime:  step.Stop,
			Logs:      convertLogs(step.Logs),
		}
	}

	links := make([]*pb.LinkRead, len(r.Links()))
	for i, link := range r.Links() {
		links[i] = &pb.LinkRead{
			Name: link.Name,
			Url:  link.URL,
		}
	}

	return &pb.RunRead{
		Uuid:      r.UUID(),
		Name:      r.Name(),
		Version:   r.RunVersion(),
		Status:    convertStatus(r.Status()),
		StartTime: r.StartTime(),
		StopTime:  r.StopTime(),
		Links:     links,
		Steps:     steps,
	}, nil
}

// Runs reads a summary of all runs for a particular deployment
func (s Service) Runs(ctx context.Context, req *pb.ReadRuns) (*pb.RunsRead, error) {
	if err := s.auth.Authorize(ctx, auth.Reader); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

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
				Name:      r.Name(),
				Version:   r.RunVersion(),
				Status:    convertStatus(r.Status()),
				StartTime: r.StartTime(),
				StopTime:  r.StopTime(),
			})
		}
	}
	return resp, nil
}

// Approve a step for a run
func (s Service) Approve(ctx context.Context, cmd *pb.ApproveStep) (*pb.StepApproved, error) {
	if err := s.auth.Authorize(ctx, auth.Reader); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}
	user := s.auth.User(ctx)

	approvals, err := s.aprvRead.All()
	if err != nil {
		logger.Error("could not read Approvals", "error", err)
		return nil, tw.ToTwirpError(err, "error loading approvals")
	}

	for _, a := range approvals {
		if a.RunUUID() == cmd.RunUuid {
			a.Complete(cmd.Approve, cmd.Override, &user, ddd.NewTimestamp())
			err = s.aprvWrite.Save(a)
			if err != nil {
				logger.Error("could not save approval", "error", err)
				return nil, tw.ToTwirpError(err, "error updating approval")
			}

			r, err := s.runRead.Load(cmd.RunUuid)
			if err != nil {
				logger.Error("could not load run", "error", err)
				return nil, tw.ToTwirpError(err, "error loading run")
			}

			// start the engine in the background
			eng, err := run.NewEngine(r, s.store, s.agent, git.DefaultAgent{}, repo.DefaultAgent{}, s.deploymentTrigger)
			if err != nil {
				logger.Error("could not recreate run engine", "error", err)
				return nil, tw.ToTwirpError(err, "error resuming run")
			}
			go eng.Resume()

			return &pb.StepApproved{}, nil
		}
	}
	return nil, twirp.NotFoundError("approval not found")
}

// Approvals gets all approvals for the user
func (s Service) Approvals(ctx context.Context, req *pb.ReadApprovals) (*pb.ApprovalsRead, error) {
	if err := s.auth.Authorize(ctx, auth.Reader); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}
	user := s.auth.User(ctx)

	approvals, err := s.aprvRead.All()
	if err != nil {
		logger.Error("could not read Approvals", "error", err)
		return nil, tw.ToTwirpError(err, "error loading approvals")
	}

	results := make([]*pb.ApprovalRead, 0)
	for _, a := range approvals {
		added := false
		for _, u := range a.Users() {
			if u == user.Name && !added {
				added = s.addApproval(a, &results)
			}
		}

		if !added {
			for _, r := range a.Roles() {
				if user.HasRole(r) && !added {
					added = s.addApproval(a, &results)
				}
			}
		}
	}

	return &pb.ApprovalsRead{
		Approvals: results,
	}, nil
}

// Manual triggers a deployment manually
func (s Service) Manual(ctx context.Context, cmd *tr.TriggerManual) (*tr.ManualTriggered, error) {
	ctx, err := s.auth.Authenticate(ctx)
	if err != nil {
		logger.Error("error authenticating for manual trigger", "error", err)
		return nil, twirp.NewError(twirp.Unauthenticated, "could not authenticate user")
	}
	user := s.auth.User(ctx)

	runUUID, err := s.trigger(&trigger.Trigger{
		UUID:      cmd.Uuid,
		Name:      cmd.Name,
		Type:      trigger.ManualTrigger,
		User:      &user,
		Arguments: []byte(cmd.Args),
	})
	if err != nil {
		return nil, tw.ToTwirpError(err, "deployment trigger failed")
	}

	return &tr.ManualTriggered{
		RunUuid: runUUID,
	}, nil
}

// Web triggers a deployment from a web call
func (s Service) Web(ctx context.Context, cmd *tr.TriggerWeb) (*tr.WebTriggered, error) {
	d, err := s.read.Load(cmd.Uuid)
	if err != nil {
		logger.Error("error reading Deployment", "error", err)
		return nil, tw.ToTwirpError(err, "couldn't read deployment to trigger")
	}

	if cmd.Token != d.Token() {
		return nil, twirp.NewError(twirp.Unauthenticated, "invalid token")
	}

	runUUID, err := s.trigger(&trigger.Trigger{
		UUID:      cmd.Uuid,
		Name:      cmd.Name,
		Type:      trigger.WebTrigger,
		Token:     &cmd.Token,
		Arguments: []byte(cmd.Args),
	})

	if err != nil {
		return nil, tw.ToTwirpError(err, "deployment trigger failed")
	}
	return &tr.WebTriggered{
		RunUuid: runUUID,
	}, nil
}

// Ready implements the ReadyService method so this service can be part of a health check routine
func (s Service) Ready() error {
	return s.ready()
}

func (s Service) addApproval(a *approval.Approval, list *[]*pb.ApprovalRead) bool {
	run, err := s.runRead.Load(a.RunUUID())
	if err != nil {
		logger.Error("could not find run", "run", a.RunUUID())
		return false
	}

	steps := run.Steps()
	stepName := steps[len(steps)-1].Name

	*list = append(*list, &pb.ApprovalRead{
		Uuid:       a.UUID(),
		RunUuid:    a.RunUUID(),
		RunName:    run.Name(),
		RunVersion: run.RunVersion(),
		StepName:   stepName,
	})
	return true
}

func (s Service) deploymentTrigger(name, deployment string, args []byte) (string, error) {
	deps, err := s.read.All()
	if err != nil {
		logger.Error("error reading Deployment", "falserror", err)
		return "", err
	}

	var d *Deployment
	for _, dep := range deps {
		if deployment == dep.Name() {
			d = dep
			break
		}
	}
	if d == nil {
		return "", errors.New("deployment not found")
	}

	return s.trigger(&trigger.Trigger{
		UUID:      d.UUID(),
		Name:      name,
		Type:      trigger.DeploymentTrigger,
		Arguments: args,
	})
}

func (s Service) trigger(trig *trigger.Trigger) (string, error) {
	// template/validate the deployment
	d, err := s.read.Load(trig.UUID)
	if err != nil {
		logger.Error("error reading Deployment", "error", err)
		return "", err
	}

	rep, err := s.repo.Load(d.RepoID())
	if err != nil {
		logger.Error("error reading Repo for Deployment", "error", err)
		return "", err
	}

	f := s.git.GetFile(rep.Endpoint(), rep.Token(), d.Branch(), d.FilePath())
	if f == nil {
		logger.Error("could not get deployment file")
		return "", err
	}

	te := template.NewEngine(&git.DefaultAgent{}, s.repo)
	temp, err := te.Run(f, trig.Arguments)
	if err != nil {
		logger.Error("yaml file could not be templated", "error", err)
		return "", err
	}
	if err = temp.Validate(); err != nil {
		logger.Error("could not validate template", "error", err)
		return "", err
	}

	if err = trig.Validate(temp); err != nil {
		logger.Error("invalid trigger", "error", err)
		return "", err
	}

	// create the run
	r, err := run.Create(ddd.NewUUID(), temp, trig)

	// start the engine in the background
	eng, err := run.NewEngine(r, s.store, s.agent, git.DefaultAgent{}, repo.DefaultAgent{}, s.deploymentTrigger)
	if err != nil {
		return "", err
	}
	go eng.Start()

	// return the run id
	return r.UUID(), nil
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
	case run.AwaitingApproval:
		return pb.Status_AWAITING_APPROVAL
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
