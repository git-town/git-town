//nolint:ireturn
package steps

import (
	"github.com/git-town/git-town/v7/src/drivers"
	"github.com/git-town/git-town/v7/src/git"
)

// SetParentBranchStep registers the branch with the given name as a parent
// of the branch with the other given name.
type SetParentBranchStep struct {
	NoOpStep
	BranchName       string
	ParentBranchName string

	previousParent string
}

func (step *SetParentBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	if step.previousParent == "" {
		return &DeleteParentBranchStep{BranchName: step.BranchName}, nil
	}
	return &SetParentBranchStep{BranchName: step.BranchName, ParentBranchName: step.previousParent}, nil
}

func (step *SetParentBranchStep) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	step.previousParent = repo.Config.ParentBranch(step.BranchName)
	return repo.Config.SetParentBranch(step.BranchName, step.ParentBranchName)
}
