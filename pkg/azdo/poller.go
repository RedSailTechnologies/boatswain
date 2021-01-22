package azdo

import (
	"context"
	"fmt"

	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/build"
	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
)

type Poller struct {
	token string
	url   string
}

func (p *Poller) GetBuilds(proj string, prID int) ([]string, []string) {
	connection := azuredevops.NewPatConnection(p.url, p.token)
	ctx := context.Background()
	client, err := build.NewClient(ctx, connection)
	if err != nil {
		return nil, nil // FIXME
	}

	branch := fmt.Sprintf("refs/pull/%d/merge", prID)
	args := build.GetBuildsArgs{
		Project:      &proj,
		ReasonFilter: &build.BuildReasonValues.PullRequest,
		StatusFilter: &build.BuildStatusValues.Completed,
		ResultFilter: &build.BuildResultValues.Succeeded,
		BranchName:   &branch,
	}

	builds, err := client.GetBuilds(ctx, args)
	if err != nil {
		return nil, nil
	}

	names := make([]string, 0)
	nums := make([]string, 0)
	for _, b := range builds.Value {
		names = append(names, *b.Definition.Name)
		nums = append(nums, *b.BuildNumber)
	}
	return names, nums
}

func (p *Poller) GetProjects() []string {
	connection := azuredevops.NewPatConnection(p.url, p.token)
	ctx := context.Background()
	client, err := core.NewClient(ctx, connection)
	if err != nil {
		return nil // FIXME
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
			return nil
		}

		for _, project := range resp.Value {
			res = append(res, *project.Name)
		}
		cont = resp.ContinuationToken
	}
	return res
}

func (p *Poller) GetPullRequests(id string) ([]string, []int) {
	connection := azuredevops.NewPatConnection(p.url, p.token)
	ctx := context.Background()
	client, err := git.NewClient(ctx, connection)

	args := git.GetPullRequestsArgs{
		RepositoryId: &id,
		SearchCriteria: &git.GitPullRequestSearchCriteria{
			Status: &git.PullRequestStatusValues.Active,
		},
	}

	resp, err := client.GetPullRequests(ctx, args)
	if err != nil {
		return nil, nil
	}

	names := make([]string, 0)
	ids := make([]int, 0)
	for _, pr := range *resp {
		names = append(names, *pr.Title)
		ids = append(ids, *pr.PullRequestId)
	}
	return names, ids
}

func (p *Poller) GetPullRequestReviewers(repoId string, prId int) []string {
	connection := azuredevops.NewPatConnection(p.url, p.token)
	ctx := context.Background()
	client, err := git.NewClient(ctx, connection)

	args := git.GetPullRequestReviewersArgs{
		RepositoryId:  &repoId,
		PullRequestId: &prId,
	}

	resp, err := client.GetPullRequestReviewers(ctx, args)
	if err != nil {
		return nil
	}

	res := make([]string, 0)
	for _, r := range *resp {
		res = append(res, *r.DisplayName)
	}
	return res
}

func (p *Poller) GetPullRequestStatuses(repoId string, prId int) []string {
	connection := azuredevops.NewPatConnection(p.url, p.token)
	ctx := context.Background()
	client, err := git.NewClient(ctx, connection)

	args := git.GetPullRequestStatusesArgs{
		RepositoryId:  &repoId,
		PullRequestId: &prId,
	}

	resp, err := client.GetPullRequestStatuses(ctx, args)
	if err != nil {
		return nil
	}

	res := make([]string, 0)
	for _, st := range *resp {
		if st.Context != nil && st.Context.Name != nil && st.Context.Genre != nil {
			res = append(res, *st.Context.Name+"/"+*st.Context.Genre)
		}
	}
	return res
}

func (p *Poller) GetRepositories(proj string) ([]string, []string) {
	connection := azuredevops.NewPatConnection(p.url, p.token)
	ctx := context.Background()
	client, err := git.NewClient(ctx, connection)

	args := git.GetRepositoriesArgs{
		Project: &proj,
	}

	resp, err := client.GetRepositories(ctx, args)
	if err != nil {
		return nil, nil
	}

	resName := make([]string, 0)
	resID := make([]string, 0)
	for _, r := range *resp {
		resName = append(resName, *r.Name)
		resID = append(resID, r.Id.String())
	}
	return resName, resID
}
