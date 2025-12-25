package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// CheckoutIfExists does the same as Checkout
// but only if that branch actually exists.
type CheckoutIfExists struct {
	Branch gitdomain.LocalBranchName
}

func (self *CheckoutIfExists) Run(args shared.RunArgs) error {
	currentBranchOpt, err := args.Git.CurrentBranch(args.Backend)
	if err != nil {
		return err
	}
	currentBranch, hasCurrentBranch := currentBranchOpt.Get()
	if !hasCurrentBranch {
		return errors.New(messages.CurrentBranchCannotDetermine)
	}
	if currentBranch == self.Branch {
		return nil
	}
	if !args.Git.BranchExists(args.Backend, self.Branch) {
		return nil
	}
	args.PrependOpcodes(&CheckoutIfNeeded{Branch: self.Branch})
	return nil
}
