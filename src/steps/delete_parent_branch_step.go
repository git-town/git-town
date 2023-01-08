package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// DeleteParentBranchStep removes the parent branch entry in the Git Town configuration.
type DeleteParentBranchStep struct {
	NoOpStep
	BranchName     string
	previousParent string
}

func (step *DeleteParentBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) { //nolint:ireturn
	if step.previousParent == "" {
		return &NoOpStep{}, nil
	}
	return &SetParentBranchStep{BranchName: step.BranchName, ParentBranchName: step.previousParent}, nil
}

func (step *DeleteParentBranchStep) Run(repo *git.ProdRepo, driver hosting.Driver) error {
	step.previousParent = repo.Config.ParentBranch(step.BranchName)
	return repo.Config.RemoveParentBranch(step.BranchName)
}
