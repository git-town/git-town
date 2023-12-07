package opcode

import (
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/vm/shared"
)

// CheckoutIfExists does the same as Checkout
// but only if that branch actually exists.
type CheckoutIfExists struct {
	Branch domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *CheckoutIfExists) Run(args shared.RunArgs) error {
	if !args.Runner.Backend.HasLocalBranch(self.Branch) {
		return nil
	}
	return (&Checkout{Branch: self.Branch}).Run(args)
}
