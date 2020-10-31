package gyarados

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"

	pb "github.com/redsailtechnologies/boatswain/rpc/gyarados"
)

func TestValuesSubstituteCorrectly(t *testing.T) {
	template := `
name: "{{ .Delivery.name }}"
helm:
  chart: "{{ .Inputs.chart.name }}"
  repo: "{{ .Parameters.repo }}"
`

	vals := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(`
Delivery:
  name: aname
Inputs:
  chart:
    name: chartname
Parameters:
  repo: myrepo
`), &vals)
	if err != nil {
		t.Error(err)
	}

	sut := &Deployment{}
	err = sut.YAML([]byte(template))
	if err != nil {
		t.Error(err)
	}

	assert.Nil(t, sut.SubsituteValues(&vals))
	assert.Equal(t, "aname", sut.Deployment.Name)
	assert.Equal(t, "chartname", sut.Deployment.Helm.Chart)
	assert.Equal(t, "myrepo", sut.Deployment.Helm.Repo)
}

func TestValidUnmarshal(t *testing.T) {
	inputs := []string{`
name: docker-deployment
docker:
  image: dockerimage
  tag: atag
`, `
name: helm-deployment
helm:
  chart: ahelmchart
  repo: myrepo
  version: 0.1.0
`, `
template: sometemplate
arguments: "{{ .Parameters.value }}"
`,
	}

	for _, str := range inputs {
		sut := &Deployment{&pb.Deployment{}, &pb.Template{}}
		err := sut.YAML([]byte(str))
		assert.Nil(t, err)
		if str[1:9] != "template" {
			assert.NotEqual(t, "", sut.Deployment.Name)
		} else {
			assert.NotEqual(t, "", sut.Template.Template)
			assert.NotEqual(t, "", sut.Template.Arguments)
		}
	}
}

func TestInvalidUnmarshal(t *testing.T) {
	inputs := []string{`
name: invalid:
`, `
name: valid
 docker:
   image: invalidSpacing
`, `
template: {{ .Template.without.quotes }}
`,
	}

	for _, str := range inputs {
		sut := &Deployment{&pb.Deployment{}, &pb.Template{}}
		err := sut.YAML([]byte(str))
		assert.NotNil(t, err)
		assert.Equal(t, "", sut.Name)
	}
}
