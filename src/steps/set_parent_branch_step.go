package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// SetParentBranchStep registers the branch with the given name as a parent
// of the branch with the other given name.
type SetParentBranchStep struct {
	NoOpStep
	BranchName       string
	ParentBranchName string
	previousParent   string
}

func (step *SetParentBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) { //nolint:ireturn
	if step.previousParent == "" {
		return &DeleteParentBranchStep{BranchName: step.BranchName}, nil
	}
	return &SetParentBranchStep{BranchName: step.BranchName, ParentBranchName: step.previousParent}, nil
}

func (step *SetParentBranchStep) Run(repo *git.ProdRepo, driver hosting.Driver) error {
	step.previousParent = repo.Config.Ancestry.Parent(step.BranchName)
	return repo.Config.Ancestry.SetParent(step.BranchName, step.ParentBranchName)
}
