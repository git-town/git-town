package opcodes

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// CheckoutParent checks out the parent branch of the current branch.
type CheckoutParent struct {
	CurrentBranch gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *CheckoutParent) Run(args shared.RunArgs) error {
	parentPtr := args.Lineage.Parent(self.CurrentBranch)
	if parentPtr == nil {
		return nil
	}
	parent := *parentPtr
	if parent == self.CurrentBranch {
		return nil
	}
	return args.Runner.Frontend.CheckoutBranch(parent, false)
}
