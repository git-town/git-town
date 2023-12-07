package opcode

import (
	"fmt"

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
	fmt.Println("0000000000000000000", currentBranch)
	self.PreviousBranchCandidates.Remove(currentBranch)

	actualPreviousBranch := args.Runner.Backend.PreviouslyCheckedOutBranch()
	expectedPreviousBranch := firstExistingBranch(self.PreviousBranchCandidates, args.Runner.Backend)
	fmt.Println("1111111111111111111", actualPreviousBranch)
	fmt.Println("2222222222222222222", expectedPreviousBranch)
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
