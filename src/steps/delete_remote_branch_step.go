package steps

import "github.com/git-town/git-town/src/git"

// DeleteRemoteBranchStep deletes the current branch from the origin remote.
type DeleteRemoteBranchStep struct {
	NoOpStep
	BranchName string
	IsTracking bool

	branchSha string
}

// CreateUndoStep returns the undo step for this step.
func (step *DeleteRemoteBranchStep) CreateUndoStep() Step {
	if step.IsTracking {
		return &CreateTrackingBranchStep{BranchName: step.BranchName}
	}
	return &CreateRemoteBranchStep{BranchName: step.BranchName, Sha: step.branchSha}
}

// Run executes this step.
func (step *DeleteRemoteBranchStep) Run(repo *git.ProdRepo) (err error) {
	if !step.IsTracking {
		trackingBranchName := repo.Silent.TrackingBranchName(step.BranchName)
		step.branchSha, err = repo.Silent.BranchSha(trackingBranchName)
		if err != nil {
			return err
		}
	}
	return repo.Logging.DeleteRemoteBranch(step.BranchName)
}
