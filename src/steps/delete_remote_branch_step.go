package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// DeleteRemoteBranchStep deletes the current branch from the origin remote.
type DeleteRemoteBranchStep struct {
	NoOpStep
	BranchName string
	IsTracking bool
	branchSha  string
}

func (step *DeleteRemoteBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) { //nolint:ireturn
	if step.IsTracking {
		return &CreateTrackingBranchStep{BranchName: step.BranchName}, nil
	}
	return &CreateRemoteBranchStep{BranchName: step.BranchName, Sha: step.branchSha}, nil
}

func (step *DeleteRemoteBranchStep) Run(repo *git.ProdRepo, driver hosting.Driver) (err error) {
	if !step.IsTracking {
		trackingBranchName := repo.Silent.TrackingBranchName(step.BranchName)
		step.branchSha, err = repo.Silent.ShaForBranch(trackingBranchName)
		if err != nil {
			return err
		}
	}
	return repo.Logging.DeleteRemoteBranch(step.BranchName)
}
