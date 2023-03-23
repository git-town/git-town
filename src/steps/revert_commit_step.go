package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// RevertCommitStep reverts the commit with the given sha.
type RevertCommitStep struct {
	EmptyStep
	Sha string
}

func (step *RevertCommitStep) Run(run *git.ProdRunner, connector hosting.Connector) error {
	return run.Frontend.RevertCommit(step.Sha)
}
