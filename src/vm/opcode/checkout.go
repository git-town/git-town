package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// Checkout checks out a new branch.
type Checkout struct {
	Branch domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (op *Checkout) Run(args shared.RunArgs) error {
	existingBranch, err := args.Runner.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	if existingBranch == op.Branch {
		return nil
	}
	return args.Runner.Frontend.CheckoutBranch(op.Branch)
}
