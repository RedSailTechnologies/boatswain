package trigger

import (
	"errors"

	"github.com/redsailtechnologies/boatswain/pkg/auth"
	"github.com/redsailtechnologies/boatswain/pkg/deployment/template"
	pb "github.com/redsailtechnologies/boatswain/rpc/deployment"
)

// Trigger is the way a deployment is started
type Trigger struct {
	Name string
	Type Type
	User auth.User
}

// Type represents the way a trigger is set off
type Type int

const (
	// WebTrigger is a http call based trigger
	WebTrigger Type = 0

	// ManualTrigger is a user calling a trigger
	ManualTrigger Type = 1

	// DeploymentTrigger is one deployment triggering another
	DeploymentTrigger Type = 2
)

// Validate takes a trigger call and validates it against a list of templated triggers
func Validate(u *auth.User, t *pb.TriggerDeployment, d template.Template) error {
	for _, trigger := range *d.Triggers {
		if t.Type == pb.TriggerDeployment_MANUAL {
			// FIXME - we need to verify groups somehow or just add them to the app?
			if trigger.Manual == nil {
				continue
			}
			for _, user := range trigger.Manual.Users {
				if u.Name == user {
					return nil
				}
			}
		} else if t.Type == pb.TriggerDeployment_WEB {
			if trigger.Web == nil {
				continue
			}
			if trigger.Web.Name == t.Name {
				return nil
			}
		}
	}
	return errors.New("invalid trigger") // FIXME - refactor for better reporting
}
