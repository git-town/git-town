package opcodes

import (
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/vm/shared"
)

// PreserveCheckoutHistory does stuff.
type PreserveCheckoutHistory struct {
	PreviousBranchCandidates gitdomain.LocalBranchNames
	undeclaredOpcodeMethods
}

func (self *PreserveCheckoutHistory) Run(args shared.RunArgs) error {
	if !args.Runner.Backend.CurrentBranchCache.Initialized() {
		// the branch cache is not initialized --> there were no branch changes --> no need to restore the branch history
		return nil
	}
	currentBranch := args.Runner.Backend.CurrentBranchCache.Value()
	actualPreviousBranch := args.Runner.Backend.CurrentBranchCache.Previous()
	// remove the current branch from the list of previous branch candidates because the current branch should never also be the previous branch
	candidatesWithoutCurrent := self.PreviousBranchCandidates.Remove(currentBranch)
	expectedPreviousBranch := args.Runner.Backend.FirstExistingBranch(candidatesWithoutCurrent, gitdomain.EmptyLocalBranchName())
	if expectedPreviousBranch.IsEmpty() || actualPreviousBranch == expectedPreviousBranch {
		return nil
	}
	err := args.Runner.Backend.CheckoutBranchUncached(expectedPreviousBranch)
	if err != nil {
		return err
	}
	return args.Runner.Backend.CheckoutBranchUncached(currentBranch)
}
