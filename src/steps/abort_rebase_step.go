package steps

import (
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// AbortRebaseStep represents aborting on ongoing merge conflict.
// This step is used in the abort scripts for Git Town commands.
type AbortRebaseStep struct {
	EmptyStep
}

func (step *AbortRebaseStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	return run.Frontend.AbortRebase()
}
