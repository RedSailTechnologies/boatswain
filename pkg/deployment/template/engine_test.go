package template

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	p, err := filepath.Abs("../../../docs/example.yaml")
	f, err := ioutil.ReadFile(p)
	sut := NewEngine(nil, f)
	if err == nil {
		getTemplates(sut.raw, make([]template, 0))
	}
}
