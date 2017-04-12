package steps

import (
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/script"
)

// CheckoutBranchStep checks out a new branch.
type CheckoutBranchStep struct {
	BranchName string
}

// CreateAbortStep returns the abort step for this step.
func (step CheckoutBranchStep) CreateAbortStep() Step {
	return NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step CheckoutBranchStep) CreateContinueStep() Step {
	return NoOpStep{}
}

// CreateUndoStep returns the undo step for this step.
func (step CheckoutBranchStep) CreateUndoStep() Step {
	return CheckoutBranchStep{BranchName: git.GetCurrentBranchName()}
}

// Run executes this step.
func (step CheckoutBranchStep) Run() error {
	if git.GetCurrentBranchName() != step.BranchName {
		return script.RunCommand("git", "checkout", step.BranchName)
	}
	return nil
}
