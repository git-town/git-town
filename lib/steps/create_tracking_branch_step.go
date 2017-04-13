package steps

import (
	"github.com/Originate/git-town/lib/script"
)

// CreateTrackingBranchStep pushes the current branch up to origin
// and marks it as tracking the current branch.
type CreateTrackingBranchStep struct {
	NoAutomaticAbortOnError
	NoUndoStepAfterRun
	BranchName string
}

// CreateAbortStep returns the abort step for this step.
func (step CreateTrackingBranchStep) CreateAbortStep() Step {
	return NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step CreateTrackingBranchStep) CreateContinueStep() Step {
	return NoOpStep{}
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step CreateTrackingBranchStep) CreateUndoStepBeforeRun() Step {
	return DeleteRemoteBranchStep{BranchName: step.BranchName}
}

// Run executes this step.
func (step CreateTrackingBranchStep) Run() error {
	return script.RunCommand("git", "push", "-u", "origin", step.BranchName)
}
