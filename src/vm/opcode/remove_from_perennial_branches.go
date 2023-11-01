package opcode

import (
	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/vm/shared"
)

// RemoveFromPerennialBranches removes the branch with the given name as a perennial branch.
type RemoveFromPerennialBranches struct {
	Branch domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *RemoveFromPerennialBranches) Run(args shared.RunArgs) error {
	return args.Runner.Config.RemoveFromPerennialBranches(self.Branch)
}
