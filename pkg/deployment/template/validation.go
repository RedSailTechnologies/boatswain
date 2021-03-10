package template

import "fmt"

var rules = []func(t *Template) []error{
	stepsMustHaveOnlyOneAction,
}

// ValidationError is an error found when validating this template
type ValidationError struct {
	m string
}

func (e ValidationError) Error() string {
	return e.m
}

// Validate checks a template for correctness so we can catch obvious error early
func (t *Template) Validate() []error {
	errs := make([]error, 0)
	for _, rule := range rules {
		err := rule(t)
		if err != nil {
			errs = append(errs, err...)
		}
	}
	return errs
}

func stepsMustHaveOnlyOneAction(t *Template) []error {
	errors := make([]error, 0)
	for i, step := range *t.Strategy {
		count := 0
		if step.Helm != nil {
			count++
		}
		if step.Approval != nil {
			count++
		}
		if step.Trigger != nil {
			count++
		}

		if count != 1 {
			err := ValidationError{
				m: fmt.Sprintf("step %d must have one, and only one, of helm, approval, or trigger defined", i),
			}
			errors = append(errors, err)
		}
	}
	return errors
}
