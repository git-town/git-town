package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// CreateRemoteBranchStep pushes the current branch up to origin.
// TODO: rename to PushBranchToRemoteStep
// TODO: what is the difference to CreateTrackingBranchStep?
type CreateRemoteBranchStep struct {
	Branch     domain.LocalBranchName
	NoPushHook bool
	Sha        domain.SHA
	EmptyStep
}

func (step *CreateRemoteBranchStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	return run.Frontend.CreateRemoteBranch(step.Sha, step.Branch, step.NoPushHook)
}
