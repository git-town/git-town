package steps

import (
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/script"
)

// CreateAndCheckoutBranchStep creates a new branch and makes it the current one.
type CreateAndCheckoutBranchStep struct {
	NoAutomaticAbortOnError
	NoUndoStep
	BranchName       string
	ParentBranchName string
}

// CreateAbortStep returns the abort step for this step.
func (step CreateAndCheckoutBranchStep) CreateAbortStep() Step {
	return NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step CreateAndCheckoutBranchStep) CreateContinueStep() Step {
	return NoOpStep{}
}

// Run executes this step.
func (step CreateAndCheckoutBranchStep) Run() error {
	git.SetParentBranch(step.BranchName, step.ParentBranchName)
	return script.RunCommand("git", "checkout", "-b", step.BranchName, step.ParentBranchName)
}
