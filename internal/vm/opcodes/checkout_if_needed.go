package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// CheckoutIfNeeded checks out a new branch.
type CheckoutIfNeeded struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CheckoutIfNeeded) Run(args shared.RunArgs) error {
	existingBranch, err := args.Git.CurrentBranch(args.Backend)
	if err != nil {
		return err
	}
	if existingBranch == self.Branch {
		return nil
	}
	args.PrependOpcodes(&Checkout{Branch: self.Branch})
	return nil
}
