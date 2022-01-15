package steps

import (
	"github.com/git-town/git-town/v7/src/drivers"
	"github.com/git-town/git-town/v7/src/git"
)

// AddToPerennialBranches adds the branch with the given name as a perennial branch.
type AddToPerennialBranches struct {
	NoOpStep
	BranchName string
}

// CreateUndoStep returns the undo step for this step.
func (step *AddToPerennialBranches) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &RemoveFromPerennialBranches{BranchName: step.BranchName}, nil
}

// Run executes this step.
func (step *AddToPerennialBranches) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	return repo.Config.AddToPerennialBranches(step.BranchName)
}
