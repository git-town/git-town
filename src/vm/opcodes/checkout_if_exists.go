package opcodes

import (
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/vm/shared"
)

// CheckoutIfExists does the same as Checkout
// but only if that branch actually exists.
type CheckoutIfExists struct {
	Branch gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *CheckoutIfExists) Run(args shared.RunArgs) error {
	existingBranch, err := args.Runner.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	if existingBranch == self.Branch {
		return nil
	}
	if !args.Runner.Backend.HasLocalBranch(self.Branch) {
		return nil
	}
	return (&Checkout{Branch: self.Branch}).Run(args)
}
