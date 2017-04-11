package steps

import (
	"github.com/Originate/git-town/lib/script"
)

// AbortRebaseBranchStep represents aborting on ongoing merge conflict.
// This step is used in the abort scripts for Git Town commands.
type AbortRebaseBranchStep struct{}

func (step AbortRebaseBranchStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step AbortRebaseBranchStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step AbortRebaseBranchStep) CreateUndoStep() Step {
	return NoOpStep{}
}

// Run executes this step.
func (step AbortRebaseBranchStep) Run() error {
	return script.RunCommand("git", "rebase", "--abort")
}
