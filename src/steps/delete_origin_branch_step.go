package steps

import (
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// DeleteOriginBranchStep deletes the current branch from the origin remote.
type DeleteOriginBranchStep struct {
	EmptyStep
	Branch     string
	IsTracking bool
	NoPushHook bool
	branchSha  string
}

func (step *DeleteOriginBranchStep) CreateUndoStep(_ *git.BackendCommands) (Step, error) {
	if step.IsTracking {
		return &CreateTrackingBranchStep{Branch: step.Branch, NoPushHook: step.NoPushHook}, nil
	}
	return &CreateRemoteBranchStep{Branch: step.Branch, Sha: step.branchSha, NoPushHook: step.NoPushHook}, nil
}

func (step *DeleteOriginBranchStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	if !step.IsTracking {
		trackingBranch := git.TrackingBranch(step.Branch)
		var err error
		step.branchSha, err = run.Backend.ShaForBranch(trackingBranch)
		if err != nil {
			return err
		}
	}
	return run.Frontend.DeleteRemoteBranch(step.Branch)
}
