package gyarados

import (
	"bytes"
	"errors"
	"html/template"
	"reflect"

	pb "github.com/redsailtechnologies/boatswain/rpc/gyarados"
	"sigs.k8s.io/yaml"
)

// Deployments is the plural for Deployment
type Deployments struct {
	Deployments []*Deployment
}

// Tests is the plural for Test
type Tests struct {
	Tests []*Deployment
}

// Deployment is what we take in a delivery and create
type Deployment struct {
	*pb.Deployment
	*pb.Template
}

func (d *Deployment) SubstituteTemplates(*[]*DeliverySpec) error {
	return nil
}

func (d *Deployment) SubsituteValues(v *map[string]interface{}) error {
	return subVals(d, v)
}

func subVals(n interface{}, v *map[string]interface{}) error {
	val := reflect.ValueOf(n)
	for key, val := range n {
		t, err := template.New("").Parse(val.(string))
		if err == nil {
			buff := new(bytes.Buffer)
			err = t.Execute(buff, v)
		}
	}
}

func (d *Deployment) Validate() error {
	return nil
}

func (d *Deployment) YAML(in []byte) error {
	// parse/set the template on its own as it won't get set otherwise
	template := pb.Template{}
	if err := yaml.UnmarshalStrict([]byte(in), &template); err == nil {
		d.Template = &template
		return nil
	}

	deployment := pb.Deployment{}
	if err := yaml.UnmarshalStrict([]byte(in), &deployment); err == nil {
		d.Deployment = &deployment
		return nil
	}
	return errors.New("could not unmarshal yaml as a deployment or a template")
}
