package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// DeleteLocalBranch deletes the branch with the given name.
type DeleteLocalBranch struct {
	Branch domain.LocalBranchName
	Force  bool
	undeclaredOpcodeMethods
}

func (op *DeleteLocalBranch) Run(args shared.RunArgs) error {
	useForce := op.Force
	if !useForce {
		parent := args.Lineage.Parent(op.Branch)
		hasUnmergedCommits, err := args.Runner.Backend.BranchHasUnmergedCommits(op.Branch, parent.Location())
		if err != nil {
			return err
		}
		if hasUnmergedCommits {
			useForce = true
		}
	}
	return args.Runner.Frontend.DeleteLocalBranch(op.Branch, useForce)
}
