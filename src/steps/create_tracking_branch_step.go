package steps

import (
	"github.com/Originate/git-town/src/script"
)

// CreateTrackingBranchStep pushes the current branch up to origin
// and marks it as tracking the current branch.
type CreateTrackingBranchStep struct {
	NoOpStep
	BranchName string
}

// AddUndoSteps adds the undo steps for this step to the undo step list
func (step *CreateTrackingBranchStep) AddUndoSteps(stepList *StepList) {
	stepList.Prepend(&DeleteRemoteBranchStep{BranchName: step.BranchName})
}

// Run executes this step.
func (step *CreateTrackingBranchStep) Run() error {
	return script.RunCommand("git", "push", "-u", "origin", step.BranchName)
}
