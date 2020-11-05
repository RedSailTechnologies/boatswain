package gyarados

import (
	"bytes"
	"errors"
	"reflect"
	"regexp"
	"text/template"

	"sigs.k8s.io/yaml"
)

// Templater is functionality common to all delivery structs
type Templater interface {
	Validate() error
}

// SubstituteTemplates takes templates and replaces them with their values
func SubstituteTemplates(t Templater, tmpl *[]Templater) (Templater, error) {
	tType := reflect.TypeOf(t).Elem()
	_, found := tType.FieldByName("Template")
	if !found {
		return t, nil
	}

	val := reflect.ValueOf(t).Elem()
	tName := val.FieldByName("Template").String()
	tArgs := val.FieldByName("Arguments").String()

	for _, template := range *tmpl {
		if reflect.TypeOf(template) == reflect.TypeOf(t) {
			if template.(*Deployment).Name == tName {
				if err := SubstituteValues(template, tArgs); err != nil {
					return nil, err
				}
				return template, nil
			}
		}
	}
	return nil, errors.New("template not found")
}

type inputs struct {
	Inputs *map[string]interface{}
}

// SubstituteValues takes values and replaces them with go templates
func SubstituteValues(t Templater, v string) error {
	ins, err := unmarshalValues(v)
	if err != nil {
		return err
	}

	vals := &inputs{
		Inputs: &ins,
	}

	dep, err := yaml.Marshal(t)
	if err != nil {
		return err
	}

	unescaped := replaceTemplateEscape(string(dep))
	tpl, err := template.New("tmpl").Parse(unescaped)
	if err != nil {
		return err
	}

	out := new(bytes.Buffer)
	tpl.Execute(out, vals)

	err = YAML(t, out.String())
	if err != nil {
		return err
	}
	return nil
}

// YAML takes a string and converts it to a deployment
func YAML(t Templater, s string) error {
	return yaml.Unmarshal([]byte(s), t)
}

func replaceTemplateEscape(s string) string {
	re := regexp.MustCompile(`([^\\])(\$)`)
	return re.ReplaceAllString(s, "$1")
}

func unmarshalValues(s string) (map[string]interface{}, error) {
	inputs := map[string]interface{}{}
	err := yaml.Unmarshal([]byte(s), &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}
