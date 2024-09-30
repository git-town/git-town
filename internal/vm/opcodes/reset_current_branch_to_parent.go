package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// ResetCurrentBranch resets all commits in the current branch.
type ResetCurrentBranchToParent struct {
	CurrentBranch           gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ResetCurrentBranchToParent) Run(args shared.RunArgs) error {
	parent, hasParent := args.Config.Config.Lineage.Parent(self.CurrentBranch).Get()
	if !hasParent {
		return nil
	}
	args.PrependOpcodes(&ResetBranch{Target: parent.BranchName()})
	return nil
}
