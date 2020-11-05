package gyarados

import (
	"testing"

	pb "github.com/redsailtechnologies/boatswain/rpc/gyarados"
	"github.com/stretchr/testify/assert"
)

func TestStepValidateInvalidYAML(t *testing.T) {
	inputs := []string{
		// name or template is requred
		`
name: astep
template: steptemplate
always:
  - deployment: adeployment
`,
		`
always:
  - deployment: adeployment
`,
		// arguments only for a template
		`
name: astep
arguments: |
  arg1: val1
always:
  - deployment: adeployment
`,
		// templates can't specify anything but arguments
		`
template: steptemplate
success:
  - deployment: adeployment
`,
		`
template: steptemplate
failure:
  - deployment: adeployment
`,
		`
template: steptemplate
any:
  - deployment: adeployment
`,
		`
template: steptemplate
always:
  - deployment: adeployment
`,
		// at least one of success, failure, any, or always
		`
name: anotherstep
hold: 10m
`,
		`
name: anotherstep
success: []
failure: []
any: []
always: []
`,
		// each action must specify deployment or test (but not both)
		`
name: astep
success:
  - helm:
      wait: true
`,
		`
name: astep
success:
  - deployment: adeployment
    test: atest
`,
		`
name: astep
failure:
  - helm:
      wait: true
`,
		`
name: astep
failure:
  - deployment: adeployment
    test: atest
`,
		`
name: astep
any:
  - helm:
      wait: true
`,
		`
name: astep
any:
  - deployment: adeployment
    test: atest
`,
		`
name: astep
always:
  - helm:
      wait: true
`,
		`
name: astep
always:
  - deployment: adeployment
    test: atest
`,
		// each action can't specify both docker and helm options
		`
name: astep
success:
  - deployment: deployment
    helm:
      wait: true
    docker:
      rm: true
`,
		`
name: astep
failure:
  - deployment: deployment
    helm:
      wait: true
    docker:
      rm: true
`,
		`
name: astep
any:
  - deployment: deployment
    helm:
      wait: true
    docker:
      rm: true
`,
		`
name: astep
always:
  - deployment: deployment
    helm:
      wait: true
    docker:
      rm: true
`,
	}

	for _, input := range inputs {
		sut := &Step{&pb.Step{}}
		YAML(sut, input)
		assert.NotNil(t, sut.Validate(), input)
	}
}

func TestStepValidateValidYAML(t *testing.T) {
	inputs := []string{
		`
name: astep
hold: 10m
success:
  - deployment: adeployment
failure:
  - test: atest
    docker:
      rm: true
any: []
always:
  - test: anothertest
`,
		`
template: sometemplate
arguments: |
  some: args
`,
	}

	for _, input := range inputs {
		sut := &Step{&pb.Step{}}
		YAML(sut, input)
		assert.Nil(t, sut.Validate(), input)
	}
}
