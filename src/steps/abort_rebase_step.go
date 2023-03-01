package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// AbortRebaseBranchStep represents aborting on ongoing merge conflict.
// This step is used in the abort scripts for Git Town commands.
type AbortRebaseBranchStep struct {
	NoOpStep
}

func (step *AbortRebaseBranchStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	return repo.Logging.AbortRebase()
}
