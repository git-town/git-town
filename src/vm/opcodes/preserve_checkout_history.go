package opcodes

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// PreserveCheckoutHistory does stuff.
type PreserveCheckoutHistory struct {
	PreviousBranchCandidates gitdomain.LocalBranchNames
	undeclaredOpcodeMethods
}

func (self *PreserveCheckoutHistory) Run(args shared.RunArgs) error {
	if !args.Backend.CurrentBranchCache.Initialized() {
		// the branch cache is not initialized --> there were no branch changes --> no need to restore the branch history
		return nil
	}
	currentBranch := args.Backend.CurrentBranchCache.Value()
	actualPreviousBranch := args.Backend.CurrentBranchCache.Previous()
	// remove the current branch from the list of previous branch candidates because the current branch should never also be the previous branch
	candidatesWithoutCurrent := self.PreviousBranchCandidates.Remove(currentBranch)
	expectedPreviousBranch := args.Backend.FirstExistingBranch(candidatesWithoutCurrent, gitdomain.EmptyLocalBranchName())
	if expectedPreviousBranch.IsEmpty() || actualPreviousBranch == expectedPreviousBranch {
		return nil
	}
	err := args.Backend.CheckoutBranchUncached(expectedPreviousBranch)
	if err != nil {
		return err
	}
	return args.Backend.CheckoutBranchUncached(currentBranch)
}
