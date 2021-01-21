package template

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	tpl "text/template"

	"gopkg.in/yaml.v3"

	"github.com/redsailtechnologies/boatswain/pkg/git"
	"github.com/redsailtechnologies/boatswain/pkg/repo"
)

// The Engine is the worker that performs all templating steps
type Engine struct {
	git  git.Agent
	repo *repo.ReadRepository
}

// NewEngine initializes the engine with required dependencies
func NewEngine(g git.Agent, r *repo.ReadRepository) *Engine {
	return &Engine{
		git:  g,
		repo: r,
	}
}

// Run performs all templating steps with the given file contents
func (e *Engine) Run(f []byte, v []byte) (*Template, error) {
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

	d := Template{}
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
				t := Substitution{}
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

func (e *Engine) replaceTemplate(y map[string]interface{}, k string, t Substitution) (map[string]interface{}, error) {
	repos, err := e.repo.All()
	if err != nil {
		return nil, err
	}
	r := findRepo(t.Repo, repos)
	if r == nil {
		return nil, err
	}

	file := e.git.GetFile(r.Endpoint(), t.Branch, t.Name, "", "")
	if file == nil {
		return nil, errors.New("file not found")
	}

	var withVals string
	if t.Arguments != nil {
		withVals, err = replaceValues(string(file), ".Inputs", *t.Arguments)
		if err != nil {
			return nil, err
		}
	} else {
		withVals = string(file)
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

func findRepo(repo string, repos []*repo.Repo) *repo.Repo {
	for _, r := range repos {
		if r.Name() == repo {
			return r
		}
	}
	return nil
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
