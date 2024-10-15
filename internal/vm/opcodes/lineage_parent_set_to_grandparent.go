package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// LineageParentSetToGrandParent sets the parent branch of the given branch to the grandparent of the given branch.
type LineageParentSetToGrandParent struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *LineageParentSetToGrandParent) Run(args shared.RunArgs) error {
	parent, hasParent := args.Config.Config.Lineage.Parent(self.Branch).Get()
	if !hasParent {
		return nil
	}
	if grandParent, hasGrandParent := args.Config.Config.Lineage.Parent(parent).Get(); hasGrandParent {
		args.PrependOpcodes(&LineageParentSet{
			Branch: self.Branch,
			Parent: grandParent,
		})
	} else {
		args.PrependOpcodes(&RemoveParent{
			Branch: self.Branch,
		})
	}
	return nil
}
