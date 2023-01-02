package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// AddToPerennialBranchesStep adds the branch with the given name as a perennial branch.
type AddToPerennialBranchesStep struct {
	NoOpStep
	BranchName string
}

func (step *AddToPerennialBranchesStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) { //nolint:ireturn
	return &RemoveFromPerennialBranchesStep{BranchName: step.BranchName}, nil
}

func (step *AddToPerennialBranchesStep) Run(repo *git.ProdRepo, driver hosting.Driver) error {
	return repo.Config.PerennialBranches.Add(step.BranchName)
}
