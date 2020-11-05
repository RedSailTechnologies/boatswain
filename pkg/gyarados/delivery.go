package gyarados

import (
	"errors"

	pb "github.com/redsailtechnologies/boatswain/rpc/gyarados"
)

// Delivery is a complete cd spec
type Delivery struct {
	*pb.Delivery
}

// Validate checks a delivery's validity
func (d Delivery) Validate() error {
	var err error

	if err = d.specOrTemplateTests(); err != nil {
		return err
	}

	if d.Name != "" {
		if err = d.specTests(); err != nil {
			return err
		}
		if err = d.clusterTests(); err != nil {
			return err
		}
		if err = d.deploymentTests(); err != nil {
			return err
		}
		if err = d.triggerTests(); err != nil {
			return err
		}
		if err = d.strategyTests(); err != nil {
			return err
		}
	} else {
		if err = d.templateTests(); err != nil {
			return err
		}
	}
	return nil
}

func (d Delivery) specOrTemplateTests() error {
	if d.Name == "" && d.Template == "" {
		return errors.New("delivery name or a template is required")
	}
	if d.Name != "" && d.Template != "" {
		return errors.New("delivery cannot be both specified and templated")
	}
	return nil
}

func (d Delivery) specTests() error {
	if d.Arguments != "" {
		return errors.New("arguments cannot be set when specifying a delivery")
	}
	if d.Version == "" {
		return errors.New("delivery version is required")
	}
	if d.Application.Name == "" || d.Application.Project == "" {
		return errors.New("delivery application requires both name and project")
	}
	return nil
}

func (d Delivery) templateTests() error {
	if d.Version != "" {
		return errors.New("template cannot specify anything but template name and arguments")
	}
	if d.Application != nil && d.Application.Name != "" {
		return errors.New("template cannot specify anything but template name and arguments")
	}
	if d.Application != nil && d.Application.Project != "" {
		return errors.New("template cannot specify anything but template name and arguments")
	}

	errs := 0
	errs += len(d.Clusters)
	errs += len(d.Deployments)
	errs += len(d.Tests)
	errs += len(d.Triggers)
	errs += len(d.Strategy)
	if errs > 0 {
		return errors.New("template cannot specify anything but template name and arguments")
	}
	return nil
}

func (d Delivery) clusterTests() error {
	if d.Clusters == nil || len(d.Clusters) == 0 {
		return errors.New("at least one cluster must be specified")
	}
	return nil
}

func (d Delivery) deploymentTests() error {
	if (d.Deployments == nil || len(d.Deployments) == 0) && (d.Tests == nil || len(d.Tests) == 0) {
		return errors.New("at least one deployment or test must be specified or templated")
	}
	for _, deployment := range d.Deployments {
		d := &Deployment{deployment}
		if err := d.Validate(); err != nil {
			return err
		}
	}
	for _, test := range d.Tests {
		t := &Deployment{test}
		if err := t.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (d Delivery) triggerTests() error {
	if d.Triggers == nil || len(d.Triggers) == 0 {
		return errors.New("at least one trigger must be specified or templated")
	}
	for _, trigger := range d.Triggers {
		t := &Trigger{trigger}
		if err := t.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (d Delivery) strategyTests() error {
	if d.Strategy == nil || len(d.Strategy) == 0 {
		return errors.New("at least one strategy step must be specified or templated")
	}
	for _, step := range d.Strategy {
		s := &Step{step}
		if err := s.Validate(); err != nil {
			return err
		}
	}
	return nil
}
