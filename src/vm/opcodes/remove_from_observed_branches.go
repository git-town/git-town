package opcodes

import (
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/vm/shared"
)

// RemoveFromObservedBranches removes the branch with the given name as an observed branch.
type RemoveFromObservedBranches struct {
	Branch gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *RemoveFromObservedBranches) Run(args shared.RunArgs) error {
	return args.Runner.Config.RemoveFromObservedBranches(self.Branch)
}
