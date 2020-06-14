package steps

import (
	"github.com/git-town/git-town/src/drivers"
	"github.com/git-town/git-town/src/git"
)

// RemoveFromPerennialBranches removes the branch with the given name as a perennial branch.
type RemoveFromPerennialBranches struct {
	NoOpStep
	BranchName string
}

// CreateUndoStep returns the undo step for this step.
func (step *RemoveFromPerennialBranches) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &AddToPerennialBranches{BranchName: step.BranchName}, nil
}

// Run executes this step.
func (step *RemoveFromPerennialBranches) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	return repo.Config.RemoveFromPerennialBranches(step.BranchName)
}
