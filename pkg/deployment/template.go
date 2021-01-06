package deployment

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func Unmarshal(d []byte) map[string]interface{} {
	out := make(map[string]interface{})
	err := yaml.Unmarshal(d, &out)
	if err != nil {
		return nil // FIXME
	}

	templates := getTemplates(out, make([]Template, 0))
	if len(templates) == 0 {
		return nil
	}
	fmt.Println(templates[0].Arguments["version"].(string))

	out = recurseYAML(out)

	return out
}

type Template struct {
	Ref       string                 `yaml:"template"`
	Repo      string                 `yaml:"repo"`
	Arguments map[string]interface{} `yaml:"arguments"`
}

func getTemplates(y map[string]interface{}, l []Template) []Template {
	for key := range y {
		if value, ok := y[key]; ok {
			if key == "template" {
				t := Template{}
				b, err := yaml.Marshal(y)
				err = yaml.Unmarshal(b, &t)

				b, err = yaml.Marshal(y["arguments"].(map[string]interface{}))
				err = yaml.Unmarshal(b, &t.Arguments)
				if err == nil { // FIXME
					return append(l, t)
				}
			}
			if recurse, ok := value.(map[string]interface{}); ok {
				l = getTemplates(recurse, l)
			} else if list, ok := value.([]interface{}); ok {
				for v := range list {
					if r, ok := list[v].(map[string]interface{}); ok {
						l = getTemplates(r, l)
					}
				}
			}
		}
	}
	return l
}

func recurseYAML(y map[string]interface{}) map[string]interface{} {
	for key := range y {
		if value, ok := y[key]; ok {
			if key == "template" {
				return templateYAML(y, key)
			}
			if recurse, ok := value.(map[string]interface{}); ok {
				recurseYAML(recurse)
			} else if list, ok := value.([]interface{}); ok {
				for v := range list {
					if r, ok := list[v].(map[string]interface{}); ok {
						list[v] = recurseYAML(r)
					}
				}
			}
		}
	}
	return y
}

func templateYAML(y map[string]interface{}, k string) map[string]interface{} {
	o := make(map[string]interface{})

	// TODO AdamP - here's where, given the templates on hand, we can sub them in
	o["template substituted"] = nil
	return o
}
