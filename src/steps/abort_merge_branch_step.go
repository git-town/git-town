package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// AbortMergeBranchStep aborts the current merge conflict.
type AbortMergeBranchStep struct {
	NoOpStep
}

func (step *AbortMergeBranchStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	return repo.Logging.AbortMerge()
}
