package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// removes the branch with the given name from the observed branches list in the Git config
type BranchesObservedRemove struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchesObservedRemove) Run(args shared.RunArgs) error {
	var err error
	if args.Config.Value.NormalConfig.ObservedBranches.Contains(self.Branch) {
		err = args.Config.Value.NormalConfig.RemoveFromObservedBranches(self.Branch)
	}
	return err
}
