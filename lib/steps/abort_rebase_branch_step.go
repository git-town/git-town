package steps

import (
	"github.com/Originate/git-town/lib/script"
)

// AbortRebaseBranchStep represents aborting on ongoing merge conflict.
// This step is used in the abort scripts for Git Town commands.
type AbortRebaseBranchStep struct {
	NoExpectedError
	NoUndoStep
}

// Run executes this step.
func (step AbortRebaseBranchStep) Run() error {
	return script.RunCommand("git", "rebase", "--abort")
}
