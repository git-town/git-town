package steps

import (
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/script"
)

// ContinueRebaseBranchStep finishes an ongoing rebase operation
// assuming all conflicts have been resolved by the user.
type ContinueRebaseBranchStep struct {
	NoOpStep
}

// CreateAbortStep returns the abort step for this step.
func (step ContinueRebaseBranchStep) CreateAbortStep() Step {
	return AbortRebaseBranchStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step ContinueRebaseBranchStep) CreateContinueStep() Step {
	return step
}

// Run executes this step.
func (step ContinueRebaseBranchStep) Run() error {
	if git.IsRebaseInProgress() {
		return script.RunCommand("git", "rebase", "--continue")
	}
	return nil
}
