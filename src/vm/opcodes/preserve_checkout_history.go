package opcodes

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// PreserveCheckoutHistory does stuff.
type PreserveCheckoutHistory struct {
	PreviousBranchCandidates gitdomain.LocalBranchNames
	undeclaredOpcodeMethods  `exhaustruct:"optional"`
}

func (self *PreserveCheckoutHistory) Run(args shared.RunArgs) error {
	if !args.Git.CurrentBranchCache.Initialized() {
		// the branch cache is not initialized --> there were no branch changes --> no need to restore the branch history
		return nil
	}
	currentBranch := args.Git.CurrentBranchCache.Value()
	actualPreviousBranch := args.Git.CurrentBranchCache.Previous()
	// remove the current branch from the list of previous branch candidates because the current branch should never also be the previous branch
	candidatesWithoutCurrent := self.PreviousBranchCandidates.Remove(currentBranch)
	expectedPreviousBranch, hasExpectedPreviousBranch := args.Git.FirstExistingBranch(args.Backend, candidatesWithoutCurrent...).Get()
	if !hasExpectedPreviousBranch || actualPreviousBranch == expectedPreviousBranch {
		return nil
	}
	err := args.Git.CheckoutBranchUncached(args.Backend, expectedPreviousBranch, false)
	if err != nil {
		return err
	}
	return args.Git.CheckoutBranchUncached(args.Backend, currentBranch, false)
}
