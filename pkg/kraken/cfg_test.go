package kraken

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToClientset(t *testing.T) {
	config := ClusterConfig{
		Name:     "cluster",
		Endpoint: "www.not.real",
		Token:    "abcdefg",
		Cert:     "notarealcert...",
	}
	sut := Config{
		[]ClusterConfig{
			config,
		},
	}

	validClientset, noErr := sut.ToClientset("cluster")
	invalidClientset, err := sut.ToClientset("notreal")

	assert.NotNil(t, validClientset)
	assert.Nil(t, noErr)
	assert.Nil(t, invalidClientset)
	assert.Equal(t, errors.New("cluster not found"), err)
}

func TestYAMLUnmarshals(t *testing.T) {
	input := []byte(`clusters:
  - name: testName
    endpoint: testEndpoint
    token: testToken
    cert: testCert
`)

	f, err := ioutil.TempFile("", "*")
	if err != nil {
		t.Fatalf("error creating temp directory")
	}
	defer os.Remove(f.Name())
	f.Write(input)
	f.Close()

	sut := &Config{}
	err = sut.YAML(fmt.Sprintf(f.Name()))

	assert.Nil(t, err)
	assert.Equal(t, &Config{
		[]ClusterConfig{
			ClusterConfig{
				Name:     "testName",
				Endpoint: "testEndpoint",
				Token:    "testToken",
				Cert:     "testCert",
			},
		},
	}, sut, "values should Unmarshal properly")
}

func TestYAMLInvalidYaml(t *testing.T) {
	input := []byte(`clusters:
  - namee: testName
    endpoint: testEndpoint
	  token: testToken
	cert: testCert
`)

	f, err := ioutil.TempFile("", "*")
	if err != nil {
		t.Fatalf("error creating temp directory")
	}
	defer os.Remove(f.Name())
	f.Write(input)
	f.Close()

	sut := &Config{}
	err = sut.YAML(fmt.Sprintf(f.Name()))

	assert.Error(t, err)
}

func TestYAMLBadFile(t *testing.T) {
	badFile := "/fakedir/doesntexist"
	sut := &Config{}
	err := sut.YAML(badFile)
	assert.Error(t, err)
}

func testGetClusterConfig(t *testing.T) {
	config := ClusterConfig{
		Name:     "cluster",
		Endpoint: "www.not.real",
		Token:    "abcdefg",
		Cert:     "notarealcert...",
	}
	sut := Config{
		[]ClusterConfig{
			config,
		},
	}

	valid, noErr := sut.getClusterConfig("cluster")
	invalid, err := sut.getClusterConfig("doesn'texist")

	assert.Equal(t, &config, valid)
	assert.Nil(t, noErr)
	assert.Nil(t, invalid)
	assert.Equal(t, errors.New("cluster not found"), err)
}
