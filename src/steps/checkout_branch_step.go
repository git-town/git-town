package steps

import (
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/script"
)

// CheckoutBranchStep checks out a new branch.
type CheckoutBranchStep struct {
	NoOpStep
	BranchName string
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step *CheckoutBranchStep) CreateUndoStepBeforeRun() Step {
	return &CheckoutBranchStep{BranchName: git.GetCurrentBranchName()}
}

// Run executes this step.
func (step *CheckoutBranchStep) Run() error {
	if git.GetCurrentBranchName() != step.BranchName {
		err := script.RunCommand("git", "checkout", step.BranchName)
		if err == nil {
			git.UpdateCurrentBranchCache(step.BranchName)
		}
		return err
	}
	return nil
}
