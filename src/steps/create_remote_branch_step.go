package steps

import (
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// CreateRemoteBranchStep pushes the current branch up to origin.
type CreateRemoteBranchStep struct {
	EmptyStep
	Branch     string
	NoPushHook bool
	Sha        git.SHA
}

func (step *CreateRemoteBranchStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	return run.Frontend.CreateRemoteBranch(step.Sha, step.Branch, step.NoPushHook)
}
