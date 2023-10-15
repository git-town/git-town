package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// RemoveFromPerennialBranches removes the branch with the given name as a perennial branch.
type RemoveFromPerennialBranches struct {
	Branch domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (op *RemoveFromPerennialBranches) Run(args shared.RunArgs) error {
	return args.Runner.Config.RemoveFromPerennialBranches(op.Branch)
}
