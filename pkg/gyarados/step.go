package gyarados

import (
	"errors"

	pb "github.com/redsailtechnologies/boatswain/rpc/gyarados"
)

// Step is the building block for a delivery strategy
type Step struct {
	*pb.Step
}

// Validate checks a step for errors
func (s Step) Validate() error {
	var err error

	if err = s.specOrTemplateTests(); err != nil {
		return err
	}

	if s.Name != "" {
		if err = s.specTests(); err != nil {
			return err
		}
		if err = s.actionTests(); err != nil {
			return err
		}
		if err = s.successTests(); err != nil {
			return err
		}
		if err = s.failureTests(); err != nil {
			return err
		}
		if err = s.anyTests(); err != nil {
			return err
		}
		if err = s.alwaysTests(); err != nil {
			return err
		}
	} else {
		if err = s.templateTests(); err != nil {
			return err
		}
	}
	return nil
}

func (s Step) specOrTemplateTests() error {
	if s.Name == "" && s.Template == "" {
		return errors.New("step name or a template is required")
	}
	if s.Name != "" && s.Template != "" {
		return errors.New("step cannot be both specified and templated")
	}
	return nil
}

func (s Step) specTests() error {
	if s.Arguments != "" {
		return errors.New("arguments cannot be set when specifying a step")
	}
	return nil
}

func (s Step) templateTests() error {
	if s.Success != nil || s.Failure != nil || s.Any != nil || s.Always != nil {
		return errors.New("template cannot specify anything but template name and arguments")
	}
	return nil
}

func (s Step) actionTests() error {
	if s.Success == nil && s.Failure == nil && s.Any == nil && s.Always == nil {
		return errors.New("at least one of success, failure, any, or always must be specified")
	}
	t := 0
	t += len(s.Success)
	t += len(s.Failure)
	t += len(s.Any)
	t += len(s.Always)
	if t == 0 {
		return errors.New("at least one action is required for each step")
	}
	return nil
}

func (s Step) successTests() error {
	for _, action := range s.Success {
		if action.Deployment == "" && action.Test == "" {
			return errors.New("success action can only specify a deployment or a test, not both")
		}
		if action.Deployment != "" && action.Test != "" {
			return errors.New("success action must specify a deployment or a test")
		}
		if action.Docker != nil && action.Helm != nil {
			return errors.New("success action cannot specify options for both docker and helm")
		}
	}
	return nil
}

func (s Step) failureTests() error {
	for _, action := range s.Failure {
		if action.Deployment == "" && action.Test == "" {
			return errors.New("failure action can only specify a deployment or a test, not both")
		}
		if action.Deployment != "" && action.Test != "" {
			return errors.New("failure action must specify a deployment or a test")
		}
		if action.Docker != nil && action.Helm != nil {
			return errors.New("failure action cannot specify options for both docker and helm")
		}
	}
	return nil
}

func (s Step) anyTests() error {
	for _, action := range s.Any {
		if action.Deployment == "" && action.Test == "" {
			return errors.New("any action can only specify a deployment or a test, not both")
		}
		if action.Deployment != "" && action.Test != "" {
			return errors.New("any action must specify a deployment or a test")
		}
		if action.Docker != nil && action.Helm != nil {
			return errors.New("any action cannot specify options for both docker and helm")
		}
	}
	return nil
}

func (s Step) alwaysTests() error {
	for _, action := range s.Always {
		if action.Deployment == "" && action.Test == "" {
			return errors.New("always action can only specify a deployment or a test, not both")
		}
		if action.Deployment != "" && action.Test != "" {
			return errors.New("always action must specify a deployment or a test")
		}
		if action.Docker != nil && action.Helm != nil {
			return errors.New("always action cannot specify options for both docker and helm")
		}
	}
	return nil
}
