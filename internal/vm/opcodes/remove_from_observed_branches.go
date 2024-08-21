package opcodes

import (
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/vm/shared"
)

// removes the branch with the given name from the observed branches list in the Git config
type RemoveFromObservedBranches struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RemoveFromObservedBranches) Run(args shared.RunArgs) error {
	return args.Config.RemoveFromObservedBranches(self.Branch)
}
