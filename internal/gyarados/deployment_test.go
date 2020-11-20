package gyarados

import (
	"testing"

	pb "github.com/redsailtechnologies/boatswain/rpc/gyarados"
	"github.com/stretchr/testify/assert"
)

func TestDeploymentValidateInvalidYAML(t *testing.T) {
	inputs := []string{
		// can't specify a deployment and make it a template
		`
name: adeployment
template: deploymenttemplate
helm:
  chart: achart
  repo: arepo
  version: "0.1.0"
`,
		// can't specify arguments without a template
		`
name: adeployment
docker:
  image: docker.io/image
  tag: latest
arguments: |
  someargs: someotherargs
  moreargs: that should fail
`,
		// must specify either a deployment name or a template
		`
helm:
  chart: somechart
  repo: somerepo
  version: "0.1.0"
`,
		// with a deployment spec must have either helm or docker
		`
name: adeployment
`,
		// a template must not contain a helm or docker section
		`
template: atemplate
helm:
  chart: somechart
  repo: somerepo
  version: 1.0
`,
		`
template: atemplate
docker:
  image: docker.io/image
  tag: latest
`,
		// cannot specify a docker and helm spec
		`
name: adeployment
docker:
  image: docker.io/image
  tag: latest
helm:
  chart: achart
  repo: arepo
  version: 0.1.0
`,
		// a helm spec must have all elements
		`
name: adeployment
helm:
  chart: achart
  repo: arepo
`,
		`
name: adeployment
helm:
  chart: achart
  version: 0.1.0
`,
		`
name: adeployment
helm:
  repo: arepo
  version: 0.1.0
`,
		// a docker spec must have all elements
		`
name: adeployment
docker:
  image: docker.io/image
`,
		`
name: adeployment
docker:
  tag: latest
`,
	}

	for _, input := range inputs {
		sut := &Deployment{&pb.Deployment{}}
		YAML(sut, input)
		assert.NotNil(t, sut.Validate(), input)
	}
}

func TestDeploymentValidateValidYAML(t *testing.T) {
	inputs := []string{
		// a normal helm deployment
		`
name: adeployment
helm:
  chart: achart
  repo: arepo
  version: "0.1.0"
`,
		// a normal docker deployment
		`
name: adeployment
docker:
  image: docker.io/tag
  tag: latest
`,
		// a normal template
		`
template: sometemplate
args: |
  a: b
  c: d
`,
	}

	for _, input := range inputs {
		sut := &Deployment{&pb.Deployment{}}
		YAML(sut, input)
		assert.Nil(t, sut.Validate(), input)
	}
}
