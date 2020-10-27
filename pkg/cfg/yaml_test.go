package cfg

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SubConfig struct {
	A string `yaml:"a,omitempty"`
	B string `yaml:"b,omitempty"`
}

type Config struct {
	Configs []SubConfig `yaml:"things,omitempty"`
}

func TestYAMLUnmarshals(t *testing.T) {
	input := []byte(`things:
  - a: aVal
    b: bVal
`)

	f, err := ioutil.TempFile("", "*")
	if err != nil {
		t.Fatalf("error creating temp directory")
	}
	defer os.Remove(f.Name())
	f.Write(input)
	f.Close()

	sut := &Config{}
	err = YAML(fmt.Sprintf(f.Name()), sut)

	assert.Nil(t, err)
	assert.Equal(t, &Config{
		[]SubConfig{
			SubConfig{
				A: "aVal",
				B: "bVal",
			},
		},
	}, sut, "values should Unmarshal properly")
}

func TestYAMLInvalidYamlTopKey(t *testing.T) {
	input := []byte(`thingss:
  - a: shouldNotExist
    b: neitherShouldThis
`)

	f, err := ioutil.TempFile("", "*")
	if err != nil {
		t.Fatalf("error creating temp directory")
	}
	defer os.Remove(f.Name())
	f.Write(input)
	f.Close()

	sut := &Config{}
	err = YAML(fmt.Sprintf(f.Name()), sut)

	assert.Nil(t, sut.Configs)
}

func TestYAMLInvalidYamlKeys(t *testing.T) {
	input := []byte(`things:
  - aa: notThisVal
    b: thisVal
`)

	f, err := ioutil.TempFile("", "*")
	if err != nil {
		t.Fatalf("error creating temp directory")
	}
	defer os.Remove(f.Name())
	f.Write(input)
	f.Close()

	sut := &Config{}
	err = YAML(fmt.Sprintf(f.Name()), sut)

	assert.Equal(t, "", sut.Configs[0].A)
	assert.Equal(t, "thisVal", sut.Configs[0].B)
}

func TestYAMLInvalidYaml(t *testing.T) {
	input := []byte(`things:
   - a: notThisVal
 b: thisVal
`)

	f, err := ioutil.TempFile("", "*")
	if err != nil {
		t.Fatalf("error creating temp directory")
	}
	defer os.Remove(f.Name())
	f.Write(input)
	f.Close()

	sut := &Config{}
	err = YAML(fmt.Sprintf(f.Name()), sut)

	assert.Error(t, err)
}

func TestYAMLBadFile(t *testing.T) {
	badFile := "/fakedir/doesntexist"
	sut := &Config{}
	err := YAML(badFile, sut)
	assert.Error(t, err)
}
