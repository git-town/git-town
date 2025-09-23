package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// CheckoutIfNeeded checks out a new branch.
type CheckoutIfNeeded struct {
	Branch gitdomain.LocalBranchName
}

func (self *CheckoutIfNeeded) Run(args shared.RunArgs) error {
	currentBranchOpt, err := args.Git.CurrentBranch(args.Backend)
	if err != nil {
		return err
	}
	currentBranch, hasCurrentBranch := currentBranchOpt.Get()
	if !hasCurrentBranch {
		return errors.New(messages.CurrentBranchCannotDetermine)
	}
	if currentBranch != self.Branch {
		args.PrependOpcodes(&Checkout{Branch: self.Branch})
	}
	return nil
}
