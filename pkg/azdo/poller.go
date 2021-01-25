package azdo

import (
	"context"
	"fmt"

	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/build"
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
