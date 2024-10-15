package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

type LineageRemoveBranch struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *LineageRemoveBranch) Run(args shared.RunArgs) error {
	parent, hasParent := args.Config.Config.Lineage.Parent(self.Branch).Get()
	children := args.Config.Config.Lineage.Children(self.Branch)
	for _, child := range children {
		if !hasParent {
			args.PrependOpcodes(&RemoveParent{Branch: child})
		} else {
			args.PrependOpcodes(&SetParent{Branch: child, Parent: parent})
		}
	}
	args.PrependOpcodes(&RemoveParent{Branch: self.Branch})
	args.Config.Config.Lineage.RemoveBranch(self.Branch)
	return nil
}
