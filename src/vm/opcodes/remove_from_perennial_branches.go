package opcodes

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// RemoveFromPerennialBranches removes the branch with the given name as a perennial branch.
type RemoveFromPerennialBranches struct {
	Branch gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *RemoveFromPerennialBranches) Run(args shared.RunArgs) error {
	return args.Config.RemoveFromPerennialBranches(self.Branch)
}
