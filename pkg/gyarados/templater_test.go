package gyarados

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	pb "github.com/redsailtechnologies/boatswain/rpc/gyarados"
)

func TestSubstituateTemplatesFound(t *testing.T) {
	inputs := []Templater{
		&Deployment{
			&pb.Deployment{
				Name: "templateA",
				Helm: &pb.Deployment_Helm{
					Chart:   "${{ .Inputs.someThing.doubleNesting.tripleNesting.tripleKey }}",
					Repo:    "${{ .Inputs.someThing.someNestedKey }}",
					Version: "${{ .Inputs.someThing.doubleNesting.doubleKey }}",
				},
			},
		},
		&Deployment{
			&pb.Deployment{
				Name: "templateB",
				Helm: &pb.Deployment_Helm{
					Chart: "anotherValue",
					Repo:  "notTheNestedValue",
				},
			},
		},
	}

	sut := &Deployment{&pb.Deployment{}}
	sut.Deployment.Template = "templateA"
	sut.Deployment.Arguments = `
someThing:
  someNestedKey: someNestedValue
  doubleNesting:
    doubleKey: doubleValue
    tripleNesting:
      tripleKey: tripleValue
`

	tmp, err := SubstituteTemplates(sut, &inputs)
	sut = tmp.(*Deployment)
	assert.Nil(t, err)
	assert.Equal(t, "templateA", sut.Name)
	assert.Equal(t, "tripleValue", sut.Helm.Chart)
	assert.Equal(t, "someNestedValue", sut.Helm.Repo)
	assert.Equal(t, "doubleValue", sut.Helm.Version)
}

func TestSubstituateTemplatesNotFound(t *testing.T) {
	inputs := []Templater{
		&Deployment{
			&pb.Deployment{
				Name: "templateA",
				Helm: &pb.Deployment_Helm{
					Chart: "${{ .Inputs.someKey }}",
					Repo:  "${{ .Inputs.someThing.someNestedKey }}",
				},
			},
		},
		&Deployment{
			&pb.Deployment{
				Name: "templateB",
				Helm: &pb.Deployment_Helm{
					Chart: "anotherValue",
					Repo:  "notTheNestedValue",
				},
			},
		},
	}

	sut := &Deployment{&pb.Deployment{}}
	sut.Deployment.Template = "templateC"
	sut.Deployment.Arguments = `
someKey: someValue
someThing:
  someNestedKey: someNestedValue
`

	tmp, err := SubstituteTemplates(sut, &inputs)
	assert.Nil(t, tmp)
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("template not found"), err)
}

func TestSubstituteValues(t *testing.T) {
	inputs := `
someKey: someValue
someThing:
  someNestedKey: someNestedValue
  irrelevantKey:
    sub: val
`
	sut := &Deployment{
		&pb.Deployment{
			Name: "templateA",
			Helm: &pb.Deployment_Helm{
				Chart: "${{ .Inputs.someKey }}",
				Repo:  "${{ .Inputs.someThing.someNestedKey }}",
			},
		},
	}
	SubstituteValues(sut, inputs)
	assert.Equal(t, "someValue", sut.Helm.Chart)
	assert.Equal(t, "someNestedValue", sut.Helm.Repo)
}

func TestValidDeploymentUnmarshal(t *testing.T) {
	inputs := []string{`
name: docker-deployment
docker:
  image: dockerimage
  tag: atag
`, `
name: helm-deployment
helm:
  chart: ${{ .Inputs }}
  repo: myrepo
  version: 0.1.0
`, `
template: sometemplate
arguments: |
  someargs: someargs
`, `
template: anothertemplate
`,
	}

	for i, str := range inputs {
		sut := &Deployment{&pb.Deployment{}}
		err := YAML(sut, str)
		assert.Nil(t, err)
		if i < 2 {
			assert.NotEqual(t, "", sut.Name)
		} else {
			assert.NotEqual(t, "", sut.Template)
		}
	}
}

func TestInvalidDeploymentUnmarshal(t *testing.T) {
	inputs := []string{`
name: invalid:
`, `
name: valid
 docker:
   image: invalidSpacing
`, `
name: {{ .Gotpl.without.dollarsign }}
`,
	}

	for _, str := range inputs {
		sut := &Deployment{&pb.Deployment{}}
		err := YAML(sut, str)
		assert.NotNil(t, err)
		assert.Equal(t, "", sut.Name)
	}
}

func TestValidStepUnmarshal(t *testing.T) {
	inputs := []string{`
name: step1
start:
  - deployment: adelivery
    helm:
      type: upgrade
`,
		`
template: sometemplate
displayName: run some deployment or something
arguments: |
  anargument: avalue
`,
	}

	for i, str := range inputs {
		sut := &Step{&pb.Step{}}
		err := YAML(sut, str)
		assert.Nil(t, err)
		if i < 1 {
			assert.NotEqual(t, "", sut.Name)
		} else {
			assert.NotEqual(t, "", sut.Template)
		}
	}
}

func TestInvalidStepUnmarshal(t *testing.T) {
	inputs := []string{`
name: invalid:
`, `
name: valid
 start:
   - deployment: invalidSpacing
`, `
name: {{ .Gotpl.without.dollarsign }}
`,
	}

	for _, str := range inputs {
		sut := &Step{&pb.Step{}}
		err := YAML(sut, str)
		assert.NotNil(t, err)
		assert.Equal(t, "", sut.Name)
	}
}

func TestValidTriggerUnmarshal(t *testing.T) {
	inputs := []string{`
name: trig
delivery:
  name: another-delivery
  trigger: trigger-name
`,
		`
name: trig
approval:
  groups: []
  users:
    - me
`,
		`
name: trig
manual:
  users: []
  groups:
    - mygroup
  params:
    - paramA
`,
		`
name: trig
web:
  name: mywebtrigger
  params:
    - paramA
`,
		`
template: sometemplate
arguments: |
  anargument: avalue
`,
	}

	for i, str := range inputs {
		sut := &Trigger{&pb.Trigger{}}
		err := YAML(sut, str)
		assert.Nil(t, err)
		if i < 4 {
			assert.NotEqual(t, "", sut.Trigger.Name)
		} else {
			assert.NotEqual(t, "", sut.Trigger.Template)
		}
	}
}

func TestInvalidTriggerUnmarshal(t *testing.T) {
	inputs := []string{`
name: invalid:
`, `
name: valid
 delivery:
   name: invalidSpacing
`, `
name: {{ .Gotpl.without.dollarsign }}
`,
	}

	for _, str := range inputs {
		sut := &Trigger{&pb.Trigger{}}
		err := YAML(sut, str)
		assert.NotNil(t, err)
		assert.Equal(t, "", sut.Name)
	}
}
