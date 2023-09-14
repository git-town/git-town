package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// CreateRemoteBranchStep pushes the given local branch up to origin.
type CreateRemoteBranchStep struct {
	Branch     domain.LocalBranchName
	NoPushHook bool
	SHA        domain.SHA
	EmptyStep
}

func (step *CreateRemoteBranchStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	return run.Frontend.CreateRemoteBranch(step.SHA, step.Branch.AtRemote(domain.OriginRemote), step.NoPushHook)
}
