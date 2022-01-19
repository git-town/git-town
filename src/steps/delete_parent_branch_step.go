//nolint:ireturn
package steps

import (
	"github.com/git-town/git-town/v7/src/drivers"
	"github.com/git-town/git-town/v7/src/git"
)

// DeleteParentBranchStep removes the parent branch entry in the Git Town configuration.
type DeleteParentBranchStep struct {
	NoOpStep
	BranchName string

	previousParent string
}

// CreateUndoStep returns the undo step for this step.
func (step *DeleteParentBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	if step.previousParent == "" {
		return &NoOpStep{}, nil
	}
	return &SetParentBranchStep{BranchName: step.BranchName, ParentBranchName: step.previousParent}, nil
}

// Run executes this step.
func (step *DeleteParentBranchStep) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	step.previousParent = repo.Config.ParentBranch(step.BranchName)
	return repo.Config.DeleteParentBranch(step.BranchName)
}
