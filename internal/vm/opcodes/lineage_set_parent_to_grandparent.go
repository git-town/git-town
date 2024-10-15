package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// LineageSetParentToGrandParent sets the given parent branch as the parent of the given branch.
// Use ChangeParent to change an existing parent.
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
	if !hasGrandParent {
		return nil
	}
	return args.Config.SetParent(self.Branch, grandParent)
}
