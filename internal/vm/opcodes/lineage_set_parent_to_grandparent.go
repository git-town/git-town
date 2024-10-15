package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// LineageSetParentToGrandParent sets the parent branch of the given branch to the grandparent of the given branch.
type LineageSetParentToGrandParent struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *LineageSetParentToGrandParent) Run(args shared.RunArgs) error {
	parent, hasParent := args.Config.Config.Lineage.Parent(self.Branch).Get()
	if !hasParent {
		return nil
	}
	grandParent, hasGrandParent := args.Config.Config.Lineage.Parent(parent).Get()
	if hasGrandParent {
		return args.Config.SetParent(self.Branch, grandParent)
	} else {
		args.Config.RemoveParent(self.Branch)
		return nil
	}
}
