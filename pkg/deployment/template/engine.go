package template

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	tpl "text/template"

	"gopkg.in/yaml.v3"

	"github.com/redsailtechnologies/boatswain/rpc/repo"
)

/*
3. clean up code
4. pass in file getting dependency
*/

// The Engine is the worker that performs all templating steps
type Engine struct {
	ctx  context.Context
	repo repo.Repo
}

// NewEngine initializes the engine with required dependencies
func NewEngine(c context.Context, r repo.Repo) *Engine {
	return &Engine{
		ctx:  c,
		repo: r,
	}
}

// Run performs all templating steps with the given file contents
func (e *Engine) Run(f []byte, v []byte) (*Deployment, error) {
	raw := make(map[string]interface{})
	err := yaml.Unmarshal(f, &raw)
	if err != nil {
		return nil, err
	}

	templated, err := e.replaceTemplates(raw)
	if err != nil {
		return nil, err
	}

	vals := make(map[string]interface{})
	err = yaml.Unmarshal(v, &vals)
	if err != nil {
		return nil, err
	}

	b, err := yaml.Marshal(templated)
	if err != nil {
		return nil, err
	}

	out, err := replaceValues(string(b), ".Parameters", vals)
	if err != nil {
		return nil, err
	}

	d := Deployment{}
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal([]byte(out), &d)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// Template performs all templating steps save the final values substitution and unmarshal
func (e *Engine) Template(f []byte) (string, error) {
	raw := make(map[string]interface{})
	err := yaml.Unmarshal(f, &raw)
	if err != nil {
		return "", err
	}

	templated, err := e.replaceTemplates(raw)
	if err != nil {
		return "", err
	}

	b, err := yaml.Marshal(templated)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (e *Engine) replaceTemplates(y map[string]interface{}) (map[string]interface{}, error) {
	for key := range y {
		if value, ok := y[key]; ok {
			if key == "template" {
				t := Template{}
				b, err := yaml.Marshal(y)
				err = yaml.Unmarshal(b, &t)
				if err != nil {
					return nil, err
				}

				if args, ok := y["arguments"]; ok {
					b, err = yaml.Marshal(args.(map[string]interface{}))
					if err != nil {
						return nil, err
					}

					err = yaml.Unmarshal(b, &t.Arguments)
					if err != nil {
						return nil, err
					}
				}
				y, err = e.replaceTemplate(y, key, t)
				if err != nil {
					return nil, err
				}
			}

			if recurse, ok := value.(map[string]interface{}); ok {
				_, err := e.replaceTemplates(recurse)
				if err != nil {
					return nil, err
				}
			} else if list, ok := value.([]interface{}); ok {
				for v := range list {
					if r, ok := list[v].(map[string]interface{}); ok {
						out, err := e.replaceTemplates(r)
						if err != nil {
							return nil, err
						}
						list[v] = out
					}
				}
			}
		}
	}
	return y, nil
}

func (e *Engine) replaceTemplate(y map[string]interface{}, k string, t Template) (map[string]interface{}, error) {
	r, err := e.repo.Find(e.ctx, &repo.FindRepo{Name: t.Repo})
	if err != nil {
		return nil, err
	}

	file, err := e.repo.File(e.ctx, &repo.ReadFile{
		RepoId:   r.Uuid,
		Branch:   t.Branch,
		FilePath: t.Name,
	})
	if err != nil {
		return nil, err
	}

	var withVals string
	if t.Arguments != nil {
		withVals, err = replaceValues(string(file.File), ".Inputs", *t.Arguments)
		if err != nil {
			return nil, err
		}
	} else {
		withVals = string(file.File)
	}

	o := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(withVals), &o)
	if err != nil {
		return nil, err
	}

	o, err = e.replaceTemplates(o)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func replaceValues(in, prefix string, vals map[string]interface{}) (string, error) {
	re := regexp.MustCompile(fmt.Sprintf("([^\\\\])\\${{ %s", prefix))
	escaped := re.ReplaceAllString(in, "$1{{ ")

	t, err := tpl.New("tpl").Parse(escaped)
	if err != nil {
		return "", err
	}

	out := new(bytes.Buffer)
	t.Execute(out, vals)
	return out.String(), nil
}
