package step

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// RemoveFromPerennialBranches removes the branch with the given name as a perennial branch.
type RemoveFromPerennialBranches struct {
	Branch domain.LocalBranchName
	Empty
}

func (step *RemoveFromPerennialBranches) Run(args RunArgs) error {
	return args.Runner.Config.RemoveFromPerennialBranches(step.Branch)
}
