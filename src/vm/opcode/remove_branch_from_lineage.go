package opcode

import (
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/vm/shared"
)

type RemoveBranchFromLineage struct {
	Branch domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *RemoveBranchFromLineage) Run(args shared.RunArgs) error {
	parent := args.Lineage.Parent(self.Branch)
	for _, child := range args.Lineage.Children(self.Branch) {
		if parent.IsEmpty() {
			args.Runner.Backend.GitTown.RemoveParent(child)
		} else {
			err := args.Runner.Backend.GitTown.SetParent(child, parent)
			if err != nil {
				return err
			}
		}
	}
	args.Runner.Backend.GitTown.RemoveParent(self.Branch)
	args.Lineage.RemoveBranch(self.Branch)
	return nil
}
