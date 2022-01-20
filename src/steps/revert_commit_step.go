package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// RevertCommitStep reverts the commit with the given sha.
type RevertCommitStep struct {
	NoOpStep
	Sha string
}

func (step *RevertCommitStep) Run(repo *git.ProdRepo, driver hosting.Driver) error {
	return repo.Logging.RevertCommit(step.Sha)
}
