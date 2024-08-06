package opcodes

import (
	"github.com/git-town/git-town/v14/internal/git/gitdomain"
	"github.com/git-town/git-town/v14/internal/vm/shared"
)

// RemoveFromObservedBranches removes the branch with the given name as an observed branch.
type RemoveFromObservedBranches struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RemoveFromObservedBranches) Run(args shared.RunArgs) error {
	return args.Config.RemoveFromObservedBranches(self.Branch)
}
