package deployment

import (
	"github.com/redsailtechnologies/boatswain/pkg/auth"
	"github.com/redsailtechnologies/boatswain/pkg/deployment/template"
	pb "github.com/redsailtechnologies/boatswain/rpc/deployment"
)

// TriggerError is returned when an error is found with a trigger
type TriggerError struct {
	Message string
}

func (e TriggerError) Error() string {
	return e.Message
}

// ValidateTrigger takes a trigger call and validates it against a list of templated triggers
func ValidateTrigger(u *auth.User, t *pb.TriggerDeployment, l []template.Trigger) error {
	for _, trigger := range l {
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
	return TriggerError{Message: "invalid trigger"}
}
