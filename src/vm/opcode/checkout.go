package opcode

import (
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/vm/shared"
)

// Checkout checks out a new branch.
type Checkout struct {
	Branch domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *Checkout) Run(args shared.RunArgs) error {
	existingBranch, err := args.Runner.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	if existingBranch == self.Branch {
		return nil
	}
	return args.Runner.Frontend.CheckoutBranch(self.Branch)
}
