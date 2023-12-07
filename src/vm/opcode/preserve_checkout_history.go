package opcode

import (
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/git"
	"github.com/git-town/git-town/v11/src/vm/shared"
)

// PreserveCheckoutHistory does stuff.
type PreserveCheckoutHistory struct {
	PreviousBranchCandidates domain.LocalBranchNames
	undeclaredOpcodeMethods
}

func (self *PreserveCheckoutHistory) Run(args shared.RunArgs) error {
	if !args.Runner.Backend.CurrentBranchCache.Initialized() {
		// the branch cache is not initialized --> there were no branch changes --> no need to restore the branch history
		return nil
	}
	currentBranch := args.Runner.Backend.CurrentBranchCache.Value()
	actualPreviousBranch := args.Runner.Backend.PreviouslyCheckedOutBranch()
	previousBranchCandidates := self.PreviousBranchCandidates.Remove(currentBranch)
	expectedPreviousBranch := firstExistingBranch(previousBranchCandidates, args.Runner.Backend)
	if expectedPreviousBranch.IsEmpty() {
		return nil
	}
	if actualPreviousBranch == expectedPreviousBranch {
		return nil
	}
	err := args.Runner.Backend.CheckoutBranchUncached(expectedPreviousBranch)
	if err != nil {
		return err
	}
	return args.Runner.Backend.CheckoutBranchUncached(currentBranch)
}

// firstExistingBranch provides the first branch in the given list that actually exists
func firstExistingBranch(candidates domain.LocalBranchNames, cmd git.BackendCommands) domain.LocalBranchName {
	for _, candidate := range candidates {
		if cmd.BranchExists(candidate) {
			return candidate
		}
	}
	return domain.EmptyLocalBranchName()
}
