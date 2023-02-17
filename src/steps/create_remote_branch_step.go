package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// CreateRemoteBranchStep pushes the current branch up to origin.
type CreateRemoteBranchStep struct {
	NoOpStep
	Branch     string
	NoPushHook bool
	Sha        string
}

func (step *CreateRemoteBranchStep) Run(repo *git.ProdRepo, driver hosting.Driver) error {
	return repo.Logging.CreateRemoteBranch(step.Sha, step.Branch, step.NoPushHook)
}
