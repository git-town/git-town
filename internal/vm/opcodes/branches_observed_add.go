package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// registers the branch with the given name as an observed branch in the Git config
type BranchesObservedAdd struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchesObservedAdd) Run(args shared.RunArgs) error {
	return args.Config.AddToObservedBranches(self.Branch)
}
