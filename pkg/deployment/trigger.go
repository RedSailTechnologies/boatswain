package deployment

import (
	"github.com/redsailtechnologies/boatswain/pkg/auth"
	pb "github.com/redsailtechnologies/boatswain/rpc/deployment"
)

// Trigger is the way a deployment is started
type Trigger struct {
	Name string
	Type TriggerType
	User auth.User
}

// TriggerType represents the way a trigger is set off
type TriggerType int

const (
	// WebTrigger is a http call based trigger
	WebTrigger TriggerType = 0

	// ManualTrigger is a user calling a trigger
	ManualTrigger TriggerType = 1

	// DeploymentTrigger is one deployment triggering another
	DeploymentTrigger TriggerType = 2
)

// TriggerError is returned when an error is found with a trigger
type TriggerError struct {
	Message string
}

func (e TriggerError) Error() string {
	return e.Message
}

// ValidateTrigger takes a trigger call and validates it against a list of templated triggers
// TODO - can we just move this into template validation?
func ValidateTrigger(u *auth.User, t *pb.TriggerDeployment, d Template) error {
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
	return TriggerError{Message: "invalid trigger"} // FIXME - refactor for better reporting
}
