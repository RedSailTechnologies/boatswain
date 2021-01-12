package run

import (
	"github.com/redsailtechnologies/boatswain/rpc/cluster"
	"github.com/redsailtechnologies/boatswain/rpc/repo"
)

// Engine performs all run steps
type Engine struct {
	cluster    cluster.Cluster
	repo       repo.Repo
	repository Repository
}

// Run starts the run by executing the templating and verification steps
func (e *Engine) Run(r Run) {
	// template deployment
	// verify deployment
	// for each step
	// start each action
	// log each action
}
