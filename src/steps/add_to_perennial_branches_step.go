package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// AddToPerennialBranchesStep adds the branch with the given name as a perennial branch.
type AddToPerennialBranchesStep struct {
	Branch domain.LocalBranchName
	EmptyStep
}

func (step *AddToPerennialBranchesStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&RemoveFromPerennialBranchesStep{Branch: step.Branch}}, nil
}

func (step *AddToPerennialBranchesStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	return run.Config.AddToPerennialBranches(step.Branch)
}
