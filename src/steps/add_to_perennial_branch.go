package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// AddToPerennialBranches adds the branch with the given name as a perennial branch.
type AddToPerennialBranches struct {
	NoOpStep
	BranchName string
}

func (step *AddToPerennialBranches) CreateUndoStep(repo *git.ProdRepo) (Step, error) { //nolint:ireturn
	return &RemoveFromPerennialBranches{BranchName: step.BranchName}, nil
}

func (step *AddToPerennialBranches) Run(repo *git.ProdRepo, driver hosting.Driver) error {
	return repo.Config.AddToPerennialBranches(step.BranchName)
}
