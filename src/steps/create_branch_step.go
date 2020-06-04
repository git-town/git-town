package steps

import (
	"github.com/git-town/git-town/src/drivers"
	"github.com/git-town/git-town/src/git"
)

// CreateBranchStep creates a new branch
// but leaves the current branch unchanged.
type CreateBranchStep struct {
	NoOpStep
	BranchName    string
	StartingPoint string
}

// CreateUndoStep returns the undo step for this step.
func (step *CreateBranchStep) CreateUndoStep() Step {
	return &DeleteLocalBranchStep{BranchName: step.BranchName}
}

// Run executes this step.
func (step *CreateBranchStep) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	return repo.Logging.CreateBranch(step.BranchName, step.StartingPoint)
}
