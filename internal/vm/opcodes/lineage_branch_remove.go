package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

type LineageBranchRemove struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *LineageBranchRemove) Run(args shared.RunArgs) error {
	parent, hasParent := args.Config.ValidatedConfig.Lineage.Parent(self.Branch).Get()
	children := args.Config.ValidatedConfig.Lineage.Children(self.Branch)
	for _, child := range children {
		if !hasParent {
			args.PrependOpcodes(&LineageParentRemove{Branch: child})
		} else {
			args.PrependOpcodes(&LineageParentSet{Branch: child, Parent: parent})
		}
	}
	args.PrependOpcodes(&LineageParentRemove{Branch: self.Branch})
	args.Config.ValidatedConfig.Lineage.RemoveBranch(self.Branch)
	return nil
}
