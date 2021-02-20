package trigger

import (
	"github.com/redsailtechnologies/boatswain/pkg/auth"
	"github.com/redsailtechnologies/boatswain/pkg/deployment/template"
)

// Trigger is the way a deployment is started
type Trigger struct {
	UUID      string
	Name      string
	Type      Type
	Token     *string
	User      *User
	Arguments []byte
}

// User is the user's identifying information attached to a trigger
type User struct {
	Name    string
	Email   string
	Roles   []auth.Role
	Subject string
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
func (t Trigger) Validate(temp *template.Template) error {
	if temp.Triggers == nil {
		return NotFoundError{}
	}
	for _, trigger := range *temp.Triggers {
		switch t.Type {
		case DeploymentTrigger:
			if t.Name == trigger.Deployment.Name {
				return t.validateDeploymentTrigger(trigger)
			}
		case ManualTrigger:
			if t.Name == trigger.Manual.Name {
				return t.validateManualTrigger(trigger)
			}
		case WebTrigger:
			if t.Name == trigger.Web.Name {
				return t.validateWebTrigger(trigger)
			}
		default:
			return NotFoundError{}
		}
	}
	return NotFoundError{}
}

// NotFoundError is thrown when the trigger in question isn't found
type NotFoundError struct{}

func (err NotFoundError) Error() string {
	return "trigger not found"
}

// Unauthorized is thrown when the user is neither in the group nor the role needed to trigger
type Unauthorized struct{}

func (err Unauthorized) Error() string {
	return "user not authorized to trigger this deployment"
}

func (t Trigger) validateDeploymentTrigger(temp template.Trigger) error {
	return nil
}

func (t Trigger) validateManualTrigger(temp template.Trigger) error {
	if t.User == nil {
		return Unauthorized{}
	}

	// first check roles as its a bit more inclusive
	role, err := auth.ToRole(temp.Manual.Role)
	if err != nil {
		return err
	}
	for _, userRole := range t.User.Roles {
		switch role {
		case auth.Reader:
			if userRole == auth.Reader {
				return nil
			}
			fallthrough
		case auth.Editor:
			if userRole == auth.Editor {
				return nil
			}
			fallthrough
		case auth.Admin:
			if userRole == auth.Admin {
				return nil
			}
		}
	}

	for _, user := range temp.Manual.Users {
		if user == t.User.Name {
			return nil
		}
	}
	return Unauthorized{}
}

func (t Trigger) validateWebTrigger(temp template.Trigger) error {
	return nil
}
