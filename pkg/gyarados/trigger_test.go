package gyarados

import (
	"testing"

	pb "github.com/redsailtechnologies/boatswain/rpc/gyarados"
	"github.com/stretchr/testify/assert"
)

func TestTriggerValidateInvalidYAML(t *testing.T) {
	inputs := []string{
		// must be delivery, web, approval, or manual
		`
name: trig
delivery:
  name: another-app
  trigger: app-trigger
approval:
  users:
    - me
`,
		`
name: trig
delivery:
  name: another-app
  trigger: app-trigger
manual:
  users:
    - me
`,
		`
name: trig
delivery:
  name: another-app
  trigger: app-trigger
web:
  name: coolwebtrigger
  params:
    - paramOne
`,
		`
name: trig
delivery:
  name: another-app
  trigger: app-trigger
approval:
  users:input
	- me
web:
  name: coolwebtrigger
  params:
    - paramOne
`,
		// delivery must have a name and a trigger
		`
name: trig
delivery:
  name: another-delivery
`,
		`
name: trig
delivery:
  trigger: atrigger
`,
		// approval/manual must have at least one user or group
		`
name: trig
approval:
  groups []
  users []
`,
		`
name: trig
manual:
  groups: []
  users: []
`,
		// name required if not trigger
		`
delivery:
  name: another-delivery
  trigger: atrigger
`,
		// name or trigger, not both
		`
name: something
template: sometemplate
`,
		// args only if trigger
		`
name: something
arguments: |
  some: args
`,
		// triggers can't have other conditions
		`
template: sometemplate
arguments: |
  some: args
manual:
  users:
    - me
`,
		// something must be set
		`
name: sometrigger
`,
	}

	for _, input := range inputs {
		sut := &Trigger{&pb.Trigger{}}
		YAML(sut, input)
		assert.NotNil(t, sut.Validate(), input)
	}
}

func TestTriggerValidateValidYAML(t *testing.T) {
	inputs := []string{
		// valid delivery
		`
name: trig
delivery:
  name: another-delivery
  trigger: trigger-name
`,
		// valid approval
		`
name: trig
approval:
  users: []
  groups:
    - mygroup
  params:
    - paramA
`,
		`
name: trig
approval:
  groups: []
  users:
    - me
`,
		// valid manual
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
manual:
  groups: []
  users:
    - me
`,
		// valid web
		`
name: trig
web:
  name: mywebtrigger
  params:
    - paramA
`,
		`
name: trig
web:
  name: mywebtrigger
`,
		// valid template
		`
template: sometemplate
arguments: |
  arg1: val1
`,
		`
template: someothertemplate
`,
	}

	for _, input := range inputs {
		sut := &Trigger{&pb.Trigger{}}
		YAML(sut, input)
		assert.Nil(t, sut.Validate(), input)
	}
}
