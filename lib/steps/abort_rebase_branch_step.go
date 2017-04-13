package steps

import (
	"github.com/Originate/git-town/lib/script"
)

// AbortRebaseBranchStep represents aborting on ongoing merge conflict.
// This step is used in the abort scripts for Git Town commands.
type AbortRebaseBranchStep struct {
	NoAutomaticAbortOnError
	NoUndoStep
}

// CreateAbortStep returns the abort step for this step.
func (step AbortRebaseBranchStep) CreateAbortStep() Step {
	return NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step AbortRebaseBranchStep) CreateContinueStep() Step {
	return NoOpStep{}
}

// Run executes this step.
func (step AbortRebaseBranchStep) Run() error {
	return script.RunCommand("git", "rebase", "--abort")
}
