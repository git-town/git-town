package opcode

import (
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/vm/shared"
)

// AddToPerennialBranches adds the branch with the given name as a perennial branch.
type AddToPerennialBranches struct {
	Branch gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *AddToPerennialBranches) Run(args shared.RunArgs) error {
	return args.Runner.GitTown.AddToPerennialBranches(self.Branch)
}
