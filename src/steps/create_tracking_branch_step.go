package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// CreateTrackingBranchStep pushes the current branch up to origin
// and marks it as tracking the current branch.
type CreateTrackingBranchStep struct {
	Branch     domain.LocalBranchName
	Remote     string
	NoPushHook bool
	EmptyStep
}

func (step *CreateTrackingBranchStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&DeleteRemoteBranchStep{Branch: step.Branch, Remote: step.Remote}}, nil
}

func (step *CreateTrackingBranchStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	return run.Frontend.PushBranch(git.PushArgs{
		Branch:     step.Branch,
		NoPushHook: step.NoPushHook,
		Remote:     step.Remote,
	})
}
