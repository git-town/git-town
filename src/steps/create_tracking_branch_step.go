package steps

import (
	"github.com/git-town/git-town/src/drivers"
	"github.com/git-town/git-town/src/git"
)

// CreateTrackingBranchStep pushes the current branch up to origin
// and marks it as tracking the current branch.
type CreateTrackingBranchStep struct {
	NoOpStep
	BranchName string
}

// CreateUndoStep returns the undo step for this step.
func (step *CreateTrackingBranchStep) CreateUndoStep() Step {
	return &DeleteRemoteBranchStep{BranchName: step.BranchName}
}

// Run executes this step.
func (step *CreateTrackingBranchStep) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	return repo.Logging.CreateTrackingBranch(step.BranchName)
}
