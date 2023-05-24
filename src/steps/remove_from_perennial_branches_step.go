package steps

import (
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// RemoveFromPerennialBranchesStep removes the branch with the given name as a perennial branch.
type RemoveFromPerennialBranchesStep struct {
	EmptyStep
	Branch string
}

func (step *RemoveFromPerennialBranchesStep) CreateUndoStep(backend *git.BackendCommands) (Step, error) {
	return &AddToPerennialBranchesStep{Branch: step.Branch}, nil
}

func (step *RemoveFromPerennialBranchesStep) Run(run *git.ProdRunner, connector hosting.Connector) error {
	return run.Config.RemoveFromPerennialBranches(step.Branch)
}
