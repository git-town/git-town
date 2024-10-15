package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// ChangeParent changes the parent of the given branch to the given parent.
// Use SetParent to set the parent if no parent existed before.
type LineageChangeParentToGrandParent struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *LineageChangeParentToGrandParent) Run(args shared.RunArgs) error {
	parent, hasParent := args.Config.Config.Lineage.Parent(self.Branch).Get()
	if hasParent {
		args.PrependOpcodes(&SetParent{
			Branch: self.Branch,
			Parent: parent,
		})
	} else {
		args.PrependOpcodes(&RemoveParent{Branch: self.Branch})
	}
	return nil
}
