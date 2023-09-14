package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
)

// RemoveFromPerennialBranchesStep removes the branch with the given name as a perennial branch.
type RemoveFromPerennialBranchesStep struct {
	Branch domain.LocalBranchName
	EmptyStep
}

func (step *RemoveFromPerennialBranchesStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&AddToPerennialBranchesStep{Branch: step.Branch}}, nil
}

func (step *RemoveFromPerennialBranchesStep) Run(args RunArgs) error {
	return args.Runner.Config.RemoveFromPerennialBranches(step.Branch)
}
