package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// AddToPerennialBranches adds the branch with the given name as a perennial branch.
type AddToPerennialBranches struct {
	Branch domain.LocalBranchName
	Empty
}

func (step *AddToPerennialBranches) Run(args RunArgs) error {
	return args.Runner.Config.AddToPerennialBranches(step.Branch)
}
