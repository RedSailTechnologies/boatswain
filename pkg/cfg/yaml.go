package cfg

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// YAML takes a relative filename and returns the config found in it
func YAML(file string, out interface{}) error {
	y, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(y, out); err != nil {
		return err
	}
	return nil
}
