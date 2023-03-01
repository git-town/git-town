package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// AddToPerennialBranchesStep adds the branch with the given name as a perennial branch.
type AddToPerennialBranchesStep struct {
	EmptyStep
	Branch string
}

func (step *AddToPerennialBranchesStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &RemoveFromPerennialBranchesStep{Branch: step.Branch}, nil
}

func (step *AddToPerennialBranchesStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	return repo.Config.AddToPerennialBranches(step.Branch)
}
