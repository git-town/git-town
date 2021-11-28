package steps

import (
	"github.com/git-town/git-town/v7/src/drivers"
	"github.com/git-town/git-town/v7/src/git"
)

// DeleteRemoteBranchStep deletes the current branch from the origin remote.
type DeleteRemoteBranchStep struct {
	NoOpStep
	BranchName string
	IsTracking bool

	branchSha string
}

// CreateUndoStep returns the undo step for this step.
func (step *DeleteRemoteBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	if step.IsTracking {
		return &CreateTrackingBranchStep{BranchName: step.BranchName}, nil
	}
	return &CreateRemoteBranchStep{BranchName: step.BranchName, Sha: step.branchSha}, nil
}

// Run executes this step.
func (step *DeleteRemoteBranchStep) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) (err error) {
	if !step.IsTracking {
		trackingBranchName := repo.Silent.TrackingBranchName(step.BranchName)
		step.branchSha, err = repo.Silent.ShaForBranch(trackingBranchName)
		if err != nil {
			return err
		}
	}
	return repo.Logging.DeleteRemoteBranch(step.BranchName)
}
