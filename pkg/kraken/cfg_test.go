package kraken

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	sut := new(ClusterList)
	err = sut.YAML(fmt.Sprintf(f.Name()))

	assert.Nil(t, err)
	assert.Equal(t, &ClusterList{
		[]Cluster{
			Cluster{
				Name:     "testName",
				Endpoint: "testEndpoint",
				Token:    "testToken",
				Cert:     "testCert",
			},
		},
	}, sut, "values should Unmarshal properly")
}

func testYAMLInvalidYaml(t *testing.T) {
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

	sut := new(ClusterList)
	err = sut.YAML(fmt.Sprintf(f.Name()))

	assert.Error(t, err)
}

func TestYAMLBadFile(t *testing.T) {
	badFile := "/fakedir/doesntexist"
	sut := new(ClusterList)
	err := sut.YAML(badFile)
	assert.Error(t, err)
}
