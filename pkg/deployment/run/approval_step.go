package run

import (
	"fmt"

	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/deployment/approval"
	"github.com/redsailtechnologies/boatswain/pkg/deployment/template"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
)

func (e *Engine) executeApprovalStep(step *template.Step) Status {
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
		appr, err = approval.Create(ddd.NewUUID(), e.run.UUID(), step.Approval.Users, step.Approval.Roles)
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

		if appr.Overridden() {
			e.run.AppendLog(fmt.Sprintf("step overriden by %s", approver.Name), Info, ddd.NewTimestamp())
			return Succeeded
		} else if appr.Approved() {
			e.run.AppendLog(fmt.Sprintf("step approved by %s", approver.Name), Info, ddd.NewTimestamp())
			return Succeeded
		}
		e.run.AppendLog(fmt.Sprintf("step rejected by %s", approver.Name), Info, ddd.NewTimestamp())
		return Failed
	}
	e.run.AppendLog("approval step not found or could not be created", Error, ddd.NewTimestamp())
	return Failed
}
