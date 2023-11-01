package opcode

import (
	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/vm/shared"
)

// CheckoutParent checks out the parent branch of the current branch.
type CheckoutParent struct {
	CurrentBranch domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *CheckoutParent) Run(args shared.RunArgs) error {
	parent := args.Lineage.Parent(self.CurrentBranch)
	if parent.IsEmpty() || parent == self.CurrentBranch {
		return nil
	}
	return args.Runner.Frontend.CheckoutBranch(parent)
}
