package opcode

import (
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/vm/shared"
)

// PreserveCheckoutHistory does stuff.
type PreserveCheckoutHistory struct {
	PreviousBranch domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *PreserveCheckoutHistory) Run(args shared.RunArgs) error {
	if !args.Runner.Backend.CurrentBranchCache.Initialized() {
		// the branch cache is not initialized --> there were no branch changes --> no need to restore the branch history
		return nil
	}
	actualPreviousBranch := args.Runner.Backend.CurrentBranchCache.Previous()
	if actualPreviousBranch == self.PreviousBranch {
		// the actually set previous branch is already the expected value --> nothing to do here
		return nil
	}
	currentBranch, err := args.Runner.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	err = args.Runner.Backend.CheckoutBranchUncached(self.PreviousBranch)
	if err != nil {
		return err
	}
	return args.Runner.Backend.CheckoutBranchUncached(currentBranch)
}
