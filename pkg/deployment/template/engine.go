package template

import (
	"gopkg.in/yaml.v3"

	"github.com/redsailtechnologies/boatswain/rpc/repo"
)

/*
3. clean up code
4. pass in file getting dependency
*/

// The Engine is the worker that performs all templating steps
type Engine struct {
	repo repo.Repo
}

// NewEngine initializes the engine with required dependencies
func NewEngine(r repo.Repo) *Engine {
	return &Engine{
		repo: r,
	}
}

// Run performs all templating steps with the given file contents
func (e *Engine) Run(f []byte) (*deployment, error) {
	raw := make(map[string]interface{})
	err := yaml.Unmarshal(f, &raw)
	if err != nil {
		return nil, err
	}
	templates, err := getTemplates(raw, make([]template, 0))
}

func getTemplates(y map[string]interface{}, l []template) ([]template, error) {
	for key := range y {
		if value, ok := y[key]; ok {
			if key == "template" {
				t := template{}
				b, err := yaml.Marshal(y)
				err = yaml.Unmarshal(b, &t)
				if err != nil {
					return nil, err
				}

				b, err = yaml.Marshal(y["arguments"].(map[string]interface{}))
				if err != nil {
					return nil, err
				}

				err = yaml.Unmarshal(b, &t.Arguments)
				if err != nil {
					return nil, err
				}
				return append(l, t), nil
			}

			var err error
			if recurse, ok := value.(map[string]interface{}); ok {
				l, err = getTemplates(recurse, l)
				if err != nil {
					return nil, err
				}
			} else if list, ok := value.([]interface{}); ok {
				for v := range list {
					if r, ok := list[v].(map[string]interface{}); ok {
						l, err = getTemplates(r, l)
						if err != nil {
							return nil, err
						}
					}
				}
			}
		}
	}
	return l, nil
}

func replaceTemplates(y map[string]interface{}) map[string]interface{} {
	for key := range y {
		if value, ok := y[key]; ok {
			if key == "template" {
				// get template
				// replace template
				return replaceTemplate(y, key)
			}
			if recurse, ok := value.(map[string]interface{}); ok {
				replaceTemplates(recurse)
			} else if list, ok := value.([]interface{}); ok {
				for v := range list {
					if r, ok := list[v].(map[string]interface{}); ok {
						list[v] = replaceTemplates(r)
					}
				}
			}
		}
	}
	return y
}

func replaceTemplate(y map[string]interface{}, k string) map[string]interface{} {
	o := make(map[string]interface{})

	// TODO AdamP -  replace values from inputs
	o["template substituted"] = nil
	return o
}
