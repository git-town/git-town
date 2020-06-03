package steps

import (
	"github.com/git-town/git-town/src/drivers"
	"github.com/git-town/git-town/src/git"
)

// RevertCommitStep reverts the commit with the given sha.
type RevertCommitStep struct {
	NoOpStep
	Sha string
}

// Run executes this step.
func (step *RevertCommitStep) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	return repo.Logging.RevertCommit(step.Sha)
}
