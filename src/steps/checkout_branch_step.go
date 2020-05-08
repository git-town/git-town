package steps

import (
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/script"
)

// CheckoutBranchStep checks out a new branch.
type CheckoutBranchStep struct {
	NoOpStep
	BranchName string

	previousBranchName string
}

// CreateUndoStep returns the undo step for this step.
func (step *CheckoutBranchStep) CreateUndoStep() Step {
	return &CheckoutBranchStep{BranchName: step.previousBranchName}
}

// Run executes this step.
func (step *CheckoutBranchStep) Run() error {
	step.previousBranchName = git.GetCurrentBranchName()
	if step.previousBranchName != step.BranchName {
		err := script.RunCommand("git", "checkout", step.BranchName)
		if err == nil {
			git.UpdateCurrentBranchCache(step.BranchName)
		}
		return err
	}
	return nil
}
