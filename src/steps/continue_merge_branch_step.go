package steps

import (
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/script"
)

// ContinueMergeBranchStep finishes an ongoing merge conflict
// assuming all conflicts have been resolved by the user.
type ContinueMergeBranchStep struct {
	NoOpStep
}

// CreateAbortStep returns the abort step for this step.
func (step ContinueMergeBranchStep) CreateAbortStep() Step {
	return NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step ContinueMergeBranchStep) CreateContinueStep() Step {
	return step
}

// Run executes this step.
func (step ContinueMergeBranchStep) Run() error {
	if git.IsMergeInProgress() {
		return script.RunCommand("git", "commit", "--no-edit")
	}
	return nil
}
