package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// AddToPerennialBranches adds the branch with the given name as a perennial branch.
type AddToPerennialBranches struct {
	Branch domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (step *AddToPerennialBranches) Run(args shared.RunArgs) error {
	return args.Runner.Config.AddToPerennialBranches(step.Branch)
}
