package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// RemoveFromPerennialBranchesStep removes the branch with the given name as a perennial branch.
type RemoveFromPerennialBranchesStep struct {
	NoOpStep
	BranchName string
}

func (step *RemoveFromPerennialBranchesStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) { //nolint:ireturn
	return &AddToPerennialBranchesStep{BranchName: step.BranchName}, nil
}

func (step *RemoveFromPerennialBranchesStep) Run(repo *git.ProdRepo, driver hosting.Driver) error {
	return repo.Config.PerennialBranches.Remove(step.BranchName)
}
