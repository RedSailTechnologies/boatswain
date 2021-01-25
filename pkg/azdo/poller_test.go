package azdo

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPoller(t *testing.T) {
	if val, ok := os.LookupEnv("TOKEN"); ok {
		sut := &DefaultAgent{
			token: val,
			url:   "https://dev.azure.com/smithtech/",
		}
		projects, _ := sut.GetProjects()
		fmt.Println("Projects:")
		for _, p := range projects {
			fmt.Printf("    %s\n", p)
		}
		fmt.Println()

		repos, _ := sut.GetRepositories("Hammerhead")
		id := ""
		fmt.Println("Repositories (Hammerhead):")
		for _, r := range repos {
			if r.Name == "Hammerhead" {
				id = r.ID
			}
			fmt.Printf("    %s\n", r)
		}
		fmt.Println()

		prs, prIDs := sut.GetPullRequests(id, "active")
		fmt.Println("Pull Requests (Hammerhead, Active):")
		for i, p := range prs {
			fmt.Printf("    %s %d\n", p, prIDs[i])
		}
		fmt.Println()

		// builds, buildNums := sut.GetBuilds("Hammerhead", prIDs[1])
		// fmt.Println("PR Builds Succeeded (Hammerhead, " + prs[1] + "):")
		// for i, b := range builds {
		// 	fmt.Printf("    %s    %s\n", b, buildNums[i])
		// }
		// fmt.Println()
	} else {
		assert.True(t, false)
	}
}
