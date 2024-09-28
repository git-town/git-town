package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

type RemoveBranchFromLineage struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RemoveBranchFromLineage) Run(args shared.RunArgs) error {
	parent, hasParent := args.Config.Config.Lineage.Parent(self.Branch).Get()
	children := args.Config.Config.Lineage.Children(self.Branch)
	for _, child := range children {
		if !hasParent {
			args.PrependOpcodes(&RemoveParent{})
		} else {
			err := args.Config.SetParent(child, parent)
			if err != nil {
				return err
			}
		}
	}
	args.Config.RemoveParent(self.Branch)
	args.Config.Config.Lineage.RemoveBranch(self.Branch)
	return nil
}
