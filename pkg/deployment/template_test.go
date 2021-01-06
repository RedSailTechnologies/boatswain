package deployment

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	p, err := filepath.Abs("../../docs/example.yaml")
	f, err := ioutil.ReadFile(p)
	if err == nil {
		Unmarshal(f)
	}
}
