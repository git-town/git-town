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
	parentPtr := args.Lineage.Parent(self.Branch)
	if parentPtr == nil {
		for _, child := range args.Lineage.Children(self.Branch) {
			args.Runner.Config.RemoveParent(child)
		}
	} else {
		parent := *parentPtr
		for _, child := range args.Lineage.Children(self.Branch) {
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
