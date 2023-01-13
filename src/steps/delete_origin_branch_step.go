package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// DeleteOriginBranchStep deletes the current branch from the origin remote.
type DeleteOriginBranchStep struct {
	NoOpStep
	BranchName string
	IsTracking bool
	NoPushHook bool
	branchSha  string
}

func (step *DeleteOriginBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) { //nolint:ireturn
	if step.IsTracking {
		return &CreateTrackingBranchStep{BranchName: step.BranchName, NoPushHook: step.NoPushHook}, nil
	}
	return &CreateRemoteBranchStep{BranchName: step.BranchName, Sha: step.branchSha, NoPushHook: step.NoPushHook}, nil
}

func (step *DeleteOriginBranchStep) Run(repo *git.ProdRepo, driver hosting.Driver) error {
	if !step.IsTracking {
		trackingBranchName := repo.Silent.TrackingBranchName(step.BranchName)
		var err error
		step.branchSha, err = repo.Silent.ShaForBranch(trackingBranchName)
		if err != nil {
			return err
		}
	}
	return repo.Logging.DeleteRemoteBranch(step.BranchName)
}
