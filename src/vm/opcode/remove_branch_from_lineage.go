package opcode

import (
	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/vm/shared"
)

type RemoveBranchFromLineage struct {
	Branch domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *RemoveBranchFromLineage) Run(args shared.RunArgs) error {
	parent := args.Lineage.Parent(self.Branch)
	for _, child := range args.Lineage.Children(self.Branch) {
		if parent.IsEmpty() {
			args.Runner.Backend.Config.RemoveParent(child)
		} else {
			err := args.Runner.Backend.Config.SetParent(child, parent)
			if err != nil {
				return err
			}
		}
	}
	args.Runner.Backend.Config.RemoveParent(self.Branch)
	args.Lineage.RemoveBranch(self.Branch)
	return nil
}
