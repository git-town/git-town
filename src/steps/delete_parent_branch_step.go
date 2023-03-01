package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// DeleteParentBranchStep removes the parent branch entry in the Git Town configuration.
type DeleteParentBranchStep struct {
	EmptyStep
	Branch         string
	previousParent string
}

func (step *DeleteParentBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	if step.previousParent == "" {
		return &EmptyStep{}, nil
	}
	return &SetParentStep{Branch: step.Branch, ParentBranch: step.previousParent}, nil
}

func (step *DeleteParentBranchStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	step.previousParent = repo.Config.ParentBranch(step.Branch)
	return repo.Config.RemoveParentBranch(step.Branch)
}
