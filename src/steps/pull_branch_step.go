package steps

import (
	"github.com/git-town/git-town/v7/src/drivers"
	"github.com/git-town/git-town/v7/src/git"
)

// PullBranchStep pulls the branch with the given name from the origin remote.
type PullBranchStep struct {
	NoOpStep
	BranchName string
}

func (step *PullBranchStep) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	return repo.Logging.Pull()
}
