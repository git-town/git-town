package steps

import (
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/script"
)

// CheckoutBranchStep checks out a new branch.
type CheckoutBranchStep struct {
	NoExpectedError
	NoUndoStepAfterRun
	BranchName string
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step CheckoutBranchStep) CreateUndoStepBeforeRun() Step {
	return CheckoutBranchStep{BranchName: git.GetCurrentBranchName()}
}

// Run executes this step.
func (step CheckoutBranchStep) Run() error {
	if git.GetCurrentBranchName() != step.BranchName {
		return script.RunCommand("git", "checkout", step.BranchName)
	}
	return nil
}
