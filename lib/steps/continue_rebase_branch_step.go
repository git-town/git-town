package steps

import (
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/script"
)

// ContinueRebaseBranchStep finishes an ongoing rebase operation
// assuming all conflicts have been resolved by the user.
type ContinueRebaseBranchStep struct{}

// CreateAbortStep returns the abort step for this step.
func (step ContinueRebaseBranchStep) CreateAbortStep() Step {
	return NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step ContinueRebaseBranchStep) CreateContinueStep() Step {
	return step
}

// CreateUndoStep returns the undo step for this step.
func (step ContinueRebaseBranchStep) CreateUndoStep() Step {
	return NoOpStep{}
}

// Run executes this step.
func (step ContinueRebaseBranchStep) Run() error {
	if git.IsRebaseInProgress() {
		return script.RunCommand("git", "rebase", "--continue")
	}
	return nil
}
