package opcodes

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

type RemoveBranchFromLineage struct {
	Branch gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *RemoveBranchFromLineage) Run(args shared.RunArgs) error {
	parent, hasParent := args.Lineage.Parent(self.Branch).Get()
	for _, child := range args.Lineage.Children(self.Branch) {
		if !hasParent {
			args.Runner.Config.RemoveParent(child)
		} else {
			err := args.Runner.Config.SetParent(child, parent)
			if err != nil {
				return err
			}
		}
	}
	args.Runner.Config.RemoveParent(self.Branch)
	args.Lineage.RemoveBranch(self.Branch)
	return nil
}
