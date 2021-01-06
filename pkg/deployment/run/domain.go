package run

import (
	"github.com/redsailtechnologies/boatswain/pkg/ddd"
)

type Run struct {
	events  []ddd.Event
	version int

	uuid  string
	start string
	end   string

	deployID   string
	clusterIDs []string
	repoIDs    []string

	// triggers
	// approval history? (maybe as part of the step?)
	// templatedYAML
	// executionHistoryOrLogs
}
