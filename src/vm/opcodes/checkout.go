package opcodes

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// Checkout checks out a new branch.
type Checkout struct {
	Branch gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *Checkout) Run(args shared.RunArgs) error {
	if self.Branch.IsEmpty() {
		return nil
	}
	existingBranch, err := args.Runner.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	if existingBranch == self.Branch {
		return nil
	}
	return args.Runner.Frontend.CheckoutBranch(self.Branch, false)
}
