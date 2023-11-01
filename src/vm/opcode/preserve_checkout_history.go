package opcode

import (
	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/vm/shared"
)

// PreserveCheckoutHistory does stuff.
type PreserveCheckoutHistory struct {
	InitialBranch                     domain.LocalBranchName
	InitialPreviouslyCheckedOutBranch domain.LocalBranchName
	MainBranch                        domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *PreserveCheckoutHistory) Run(args shared.RunArgs) error {
	expectedPreviouslyCheckedOutBranch, err := args.Runner.Backend.ExpectedPreviouslyCheckedOutBranch(self.InitialPreviouslyCheckedOutBranch, self.InitialBranch, self.MainBranch)
	if err != nil {
		return err
	}
	if expectedPreviouslyCheckedOutBranch == args.Runner.Backend.PreviouslyCheckedOutBranch() {
		return nil
	}
	currentBranch, err := args.Runner.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	err = args.Runner.Backend.CheckoutBranchUncached(expectedPreviouslyCheckedOutBranch)
	if err != nil {
		return err
	}
	return args.Runner.Backend.CheckoutBranchUncached(currentBranch)
}
