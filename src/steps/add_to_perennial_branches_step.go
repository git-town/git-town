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

func (step *AddToPerennialBranchesStep) CreateUndoStep(backend *git.BackendCommands) (Step, error) {
	return &RemoveFromPerennialBranchesStep{Branch: step.Branch}, nil
}

func (step *AddToPerennialBranchesStep) Run(run *git.ProdRunner, connector hosting.Connector) error {
	return run.Config.AddToPerennialBranches(step.Branch)
}
