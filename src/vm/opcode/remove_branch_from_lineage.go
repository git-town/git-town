package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

type RemoveBranchFromLineage struct {
	Branch domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (op *RemoveBranchFromLineage) Run(args shared.RunArgs) error {
	parent := args.Lineage.Parent(op.Branch)
	for _, child := range args.Lineage.Children(op.Branch) {
		if parent.IsEmpty() {
			args.Runner.Backend.Config.RemoveParent(child)
		} else {
			err := args.Runner.Backend.Config.SetParent(child, parent)
			if err != nil {
				return err
			}
		}
	}
	args.Runner.Backend.Config.RemoveParent(op.Branch)
	args.Lineage.RemoveBranch(op.Branch)
	return nil
}
