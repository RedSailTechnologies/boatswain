package gyarados

import (
	"errors"

	pb "github.com/redsailtechnologies/boatswain/rpc/gyarados"
)

// Trigger is a spec that helps determine when a delivery starts
type Trigger struct {
	*pb.Trigger
}

// Validate checks to see if the Trigger is valid
func (t Trigger) Validate() error {
	var err error

	if err = t.specOrTemplateTests(); err != nil {
		return err
	}

	if t.Name != "" {
		if err = t.specTests(); err != nil {
			return err
		}
		if err = t.singletonTests(); err != nil {
			return err
		}
		if t.Delivery != nil {
			if err = t.deliveryTests(); err != nil {
				return err
			}
		}
		if t.Approval != nil {
			if err = t.approvalTests(); err != nil {
				return err
			}
		}
		if t.Manual != nil {
			if err = t.manualTests(); err != nil {
				return err
			}
		}
	} else {
		if err = t.templateTests(); err != nil {
			return err
		}
	}
	return nil
}

func (t Trigger) specOrTemplateTests() error {
	if t.Name == "" && t.Template == "" {
		return errors.New("trigger name or a template is required")
	}
	if t.Name != "" && t.Template != "" {
		return errors.New("trigger cannot be both specified and templated")
	}
	return nil
}

func (t Trigger) specTests() error {
	if t.Arguments != "" {
		return errors.New("arguments cannot be set when specifying a template")
	}
	return nil
}

func (t Trigger) templateTests() error {
	any := t.Delivery != nil ||
		t.Approval != nil ||
		t.Manual != nil ||
		t.Web != nil
	if any {
		return errors.New("template cannot specify anything but template name and arguments")
	}
	return nil
}

func (t Trigger) singletonTests() error {
	n := 0
	if t.Delivery != nil {
		n++
	}
	if t.Approval != nil {
		n++
	}
	if t.Manual != nil {
		n++
	}
	if t.Web != nil {
		n++
	}

	if n < 1 {
		return errors.New("at least one must be specified among delivery, approval, manual, web")
	}
	if n > 1 {
		return errors.New("no more than one can be specified among or, and, delivery, approval, manual, web")
	}
	return nil
}

func (t Trigger) deliveryTests() error {
	if t.Delivery.Name == "" || t.Delivery.Trigger == "" {
		return errors.New("both delivery name and trigger must be specified")
	}
	return nil
}

func (t Trigger) approvalTests() error {
	if len(t.Approval.Groups) == 0 && len(t.Approval.Users) == 0 {
		return errors.New("approvals must specify at least one user or group")
	}
	return nil
}

func (t Trigger) manualTests() error {
	if len(t.Manual.Groups) == 0 && len(t.Manual.Users) == 0 {
		return errors.New("manual must specify at least one user or group")
	}
	return nil
}
