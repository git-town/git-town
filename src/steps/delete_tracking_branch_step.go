package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// DeleteTrackingBranchStep deletes the tracking branch of the given local branch.
type DeleteTrackingBranchStep struct {
	Branch     domain.LocalBranchName
	Remote     domain.Remote
	NoPushHook bool
	EmptyStep
}

func (step *DeleteTrackingBranchStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&CreateTrackingBranchStep{Branch: step.Branch, NoPushHook: step.NoPushHook}}, nil
}

func (step *DeleteTrackingBranchStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	return run.Frontend.DeleteRemoteBranch(step.Branch, step.Remote)
}
