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
	parent, hasParent := args.Config.Value.NormalConfig.Lineage.Parent(self.Branch).Get()
	children := args.Config.Value.NormalConfig.Lineage.Children(self.Branch)
	for _, child := range children {
		if !hasParent {
			args.PrependOpcodes(&LineageParentRemove{Branch: child})
		} else {
			args.PrependOpcodes(&LineageParentSet{Branch: child, Parent: parent})
		}
	}
	args.PrependOpcodes(&LineageParentRemove{Branch: self.Branch})
	args.Config.Value.NormalConfig.Lineage = args.Config.Value.NormalConfig.Lineage.RemoveBranch(self.Branch)
	return nil
}
