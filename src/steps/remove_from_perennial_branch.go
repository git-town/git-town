package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// RemoveFromPerennialBranches removes the branch with the given name as a perennial branch.
type RemoveFromPerennialBranches struct {
	NoOpStep
	BranchName string
}

func (step *RemoveFromPerennialBranches) CreateUndoStep(repo *git.ProdRepo) (Step, error) { //nolint:ireturn
	return &AddToPerennialBranches{BranchName: step.BranchName}, nil
}

func (step *RemoveFromPerennialBranches) Run(repo *git.ProdRepo, driver hosting.CodeHostingDriver) error {
	return repo.Config.RemoveFromPerennialBranches(step.BranchName)
}
