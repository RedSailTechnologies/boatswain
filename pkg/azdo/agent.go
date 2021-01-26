package azdo

import (
	"context"
	"errors"
	"time"

	"github.com/redsailtechnologies/boatswain/pkg/logger"

	"github.com/google/uuid"

	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/microsoft/azure-devops-go-api/azuredevops/servicehooks"
)

// Agent is the interface for interacting with an azdo project
type Agent interface {
	GetProjects() []string
	GetPullRequests(repoID, status string) ([]*PullRequest, error)
	GetRepositories(proj string) ([]*GitRepo, error)
}

// DefaultAgent is the default implementation of the azdo agent
type DefaultAgent struct {
	token string
	url   string
}

// NewDefaultAgent takes basic azdo auth info and returns a new agent
func NewDefaultAgent(u, t string) *DefaultAgent {
	return &DefaultAgent{
		token: t,
		url:   u,
	}
}

func (a *DefaultAgent) GetEvents(subID string) error {
	connection := azuredevops.NewPatConnection(a.url, a.token)
	ctx := context.Background()
	client := servicehooks.NewClient(ctx, connection)

	tempUUID, _ := uuid.Parse(subID)
	sub, _ := client.GetSubscription(ctx, servicehooks.GetSubscriptionArgs{
		SubscriptionId: &tempUUID,
	})
	if sub == nil {
		return errors.New("something")
	}

	var notifications *[]servicehooks.Notification
	for notifications == nil || len(*notifications) < 1 {
		notifications, _ = client.GetNotifications(ctx, servicehooks.GetNotificationsArgs{
			SubscriptionId: &tempUUID,
		})
		if notifications == nil {
			return errors.New("something")
		}
		logger.Info("notifications found", "len", len(*notifications))
		time.Sleep(5 * time.Millisecond)
	}

	return nil
}

// GetPullRequests gets all pull requests based on the repo id and the pr status
func (a *DefaultAgent) GetPullRequests(repoID, status string) ([]*PullRequest, error) {
	connection := azuredevops.NewPatConnection(a.url, a.token)
	ctx := context.Background()
	client, err := git.NewClient(ctx, connection)
	if err != nil {
		return nil, err
	}

	args := git.GetPullRequestsArgs{
		RepositoryId: &repoID,
		SearchCriteria: &git.GitPullRequestSearchCriteria{
			Status: stringToPRStatus(status),
		},
	}

	resp, err := client.GetPullRequests(ctx, args)
	if err != nil {
		return nil, err
	}

	res := make([]*PullRequest, 0)
	for _, pr := range *resp {
		res = append(res, &PullRequest{
			ID:    *pr.PullRequestId,
			Title: *pr.Title,
		})
	}
	return res, nil
}

// GetProjects gets the names of all projects in the org
func (a *DefaultAgent) GetProjects() ([]string, error) {
	connection := azuredevops.NewPatConnection(a.url, a.token)
	ctx := context.Background()
	client, err := core.NewClient(ctx, connection)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0)
	cont := ""
	for ok := true; ok; ok = (cont != "") {
		args := core.GetProjectsArgs{}
		if cont != "" {
			args.ContinuationToken = &cont
		}

		resp, err := client.GetProjects(ctx, args)
		if err != nil {
			return nil, err
		}

		for _, project := range resp.Value {
			res = append(res, *project.Name)
		}
		cont = resp.ContinuationToken
	}
	return res, nil
}

// GetRepositories gets all git repos given a project name
func (a *DefaultAgent) GetRepositories(proj string) ([]*GitRepo, error) {
	connection := azuredevops.NewPatConnection(a.url, a.token)
	ctx := context.Background()
	client, err := git.NewClient(ctx, connection)
	if err != nil {
		return nil, err
	}

	args := git.GetRepositoriesArgs{
		Project: &proj,
	}

	resp, err := client.GetRepositories(ctx, args)
	if err != nil {
		return nil, err
	}

	res := make([]*GitRepo, 0)
	for _, r := range *resp {
		res = append(res, &GitRepo{
			ID:   r.Id.String(),
			Name: *r.Name,
		})
	}
	return res, nil
}

func stringToPRStatus(s string) *git.PullRequestStatus {
	switch s {
	case "notSet":
		return &git.PullRequestStatusValues.NotSet
	case "active":
		return &git.PullRequestStatusValues.Active
	case "abandoned":
		return &git.PullRequestStatusValues.Abandoned
	case "completed":
		return &git.PullRequestStatusValues.Completed
	case "all":
		return &git.PullRequestStatusValues.All
	default:
		return nil
	}
}
