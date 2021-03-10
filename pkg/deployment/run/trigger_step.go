package run

import (
	"fmt"

	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/deployment/template"
	"gopkg.in/yaml.v2"
)

func (e *Engine) executeTriggerStep(step *template.Step) Status {
	b, err := yaml.Marshal(step.Trigger.Arguments)
	if err != nil {
		e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
		return Failed
	}

	res, err := e.trigger(step.Trigger.Name, step.Trigger.Deployment, b)
	if err != nil {
		e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
		return Failed
	}
	e.run.AppendLog(fmt.Sprintf("triggered run with uuid %s", res), Info, ddd.NewTimestamp())
	return Succeeded
}
