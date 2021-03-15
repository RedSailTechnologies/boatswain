package run

import (
	"fmt"

	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/deployment/approval"
	"github.com/redsailtechnologies/boatswain/pkg/deployment/template"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
)

func (e *Engine) executeApprovalStep(step *template.Step) Status {
	if step.Approval.Action == "create" {
		return e.executeApprovalCreateStep(step)
	} else {
		return e.executeApprovalCompleteStep(step)
	}
}

func (e *Engine) executeApprovalCompleteStep(step *template.Step) Status {
	approvals, err := e.aprvRead.All()
	if err != nil {
		e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
		return Failed
	}

	status := Succeeded

	for _, a := range approvals {
		var err error
		if a.Name() == step.Approval.Name {
			if step.Approval.Action == "approve" {
				err = e.approve(step.Approval.Name, e.run.Name(), true)
			} else if step.Approval.Action == "reject" {
				err = e.approve(step.Approval.Name, e.run.Name(), false)
			} else {
				e.run.AppendLog("could not parse approval action", Error, ddd.NewTimestamp())
				return Failed
			}

			if err != nil {
				e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
				status = Failed
			}
		}
	}
	return status
}

func (e *Engine) executeApprovalCreateStep(step *template.Step) Status {
	approvals, err := e.aprvRead.All()
	if err != nil {
		e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
		return Failed
	}

	var appr *approval.Approval
	for _, a := range approvals {
		if a.RunUUID() == e.run.UUID() {
			appr = a
		}
	}
	if appr == nil {
		appr, err = approval.Create(ddd.NewUUID(), e.run.UUID(), step.Approval.Name, step.Approval.Users, step.Approval.Roles)
		if err != nil {
			e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
			return Failed
		}

		err = e.aprvWrite.Save(appr)
		if err != nil {
			e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
			return Failed
		}
		e.run.AppendLog("awaiting approval", Info, ddd.NewTimestamp())
		return AwaitingApproval
	} else if appr.Completed() {
		appr.Destroy()
		err := e.aprvWrite.Save(appr)
		if err != nil {
			logger.Warn("step could not be approved", "error", err, "run", e.run.UUID())
			return AwaitingApproval
		}

		approver, err := appr.Approver()
		if err != nil {
			e.run.AppendLog(err.Error(), Error, ddd.NewTimestamp())
			return Failed
		}

		if appr.Approved() {
			e.run.AppendLog(fmt.Sprintf("step approved by %s", approver.Name), Info, ddd.NewTimestamp())
			return Succeeded
		}
		e.run.AppendLog(fmt.Sprintf("step rejected by %s", approver.Name), Info, ddd.NewTimestamp())
		return Failed
	}
	e.run.AppendLog("approval step not found or could not be created", Error, ddd.NewTimestamp())
	return Failed
}
