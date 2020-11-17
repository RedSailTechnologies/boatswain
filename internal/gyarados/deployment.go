package gyarados

import (
	"errors"

	pb "github.com/redsailtechnologies/boatswain/rpc/gyarados"
)

// Deployment is a metadata object for what is created in a deployment
type Deployment struct {
	*pb.Deployment
}

// Validate checks to see if the deployment is valid
func (d Deployment) Validate() error {
	var err error

	if err = d.specOrTemplateTests(); err != nil {
		return err
	}

	if d.Name != "" {
		if err = d.specTests(); err != nil {
			return err
		}
		if err = d.dockerOrHelmTests(); err != nil {
			return err
		}
		if err = d.dockerTests(); err != nil {
			return err
		}
		if err = d.helmTests(); err != nil {
			return err
		}
	} else {
		if err = d.templateTests(); err != nil {
			return err
		}
	}
	return nil
}

func (d Deployment) specOrTemplateTests() error {
	if d.Name == "" && d.Template == "" {
		return errors.New("deployment name or a template is required")
	}
	if d.Name != "" && d.Template != "" {
		return errors.New("deployment cannot be both specified and templated")
	}
	return nil
}

func (d Deployment) specTests() error {
	if d.Arguments != "" {
		return errors.New("arguments cannot be set when specifying a deployment")
	}
	return nil
}

func (d Deployment) templateTests() error {
	if d.Docker != nil || d.Helm != nil {
		return errors.New("template cannot specify anything but template name and arguments")
	}
	return nil
}

func (d Deployment) dockerOrHelmTests() error {
	if d.Docker != nil && d.Helm != nil {
		return errors.New("a deployment must be either helm or docker, not both")
	}
	if d.Docker == nil && d.Helm == nil {
		return errors.New("a deployment must be helm or docker, one is required")
	}
	return nil
}

func (d Deployment) dockerTests() error {
	if d.Docker != nil {
		if d.Docker.Image == "" || d.Docker.Tag == "" {
			return errors.New("docker image and tag must be specified")
		}
	}
	return nil
}

func (d Deployment) helmTests() error {
	if d.Helm != nil {
		if d.Helm.Chart == "" || d.Helm.Repo == "" || d.Helm.Version == "" {
			return errors.New("helm chart, repo, and version must all be specified")
		}
	}
	return nil
}
