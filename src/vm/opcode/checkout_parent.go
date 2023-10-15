package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// CheckoutParent checks out the parent branch of the current branch.
type CheckoutParent struct {
	CurrentBranch domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (step *CheckoutParent) Run(args shared.RunArgs) error {
	parent := args.Lineage.Parent(step.CurrentBranch)
	if parent.IsEmpty() || parent == step.CurrentBranch {
		return nil
	}
	return args.Runner.Frontend.CheckoutBranch(parent)
}
