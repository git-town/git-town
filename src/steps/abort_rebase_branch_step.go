package steps

import (
	"github.com/git-town/git-town/src/git"
)

// AbortRebaseBranchStep represents aborting on ongoing merge conflict.
// This step is used in the abort scripts for Git Town commands.
type AbortRebaseBranchStep struct {
	NoOpStep
}

// Run executes this step.
func (step *AbortRebaseBranchStep) Run(repo *git.ProdRepo) error {
	return repo.Logging.AbortRebase()
}
