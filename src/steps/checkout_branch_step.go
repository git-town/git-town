package steps

import (
	"github.com/git-town/git-town/src/git"
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
func (step *CheckoutBranchStep) Run(runner *git.Runner) (err error) {
	step.previousBranchName, err = runner.CurrentBranch()
	if err != nil {
		return err
	}
	if step.previousBranchName != step.BranchName {
		return runner.CheckoutBranch(step.BranchName)
	}
	return nil
}
