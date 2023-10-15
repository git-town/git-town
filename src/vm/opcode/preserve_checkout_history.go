package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// PreserveCheckoutHistory does stuff.
type PreserveCheckoutHistory struct {
	InitialBranch                     domain.LocalBranchName
	InitialPreviouslyCheckedOutBranch domain.LocalBranchName
	MainBranch                        domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (step *PreserveCheckoutHistory) Run(args shared.RunArgs) error {
	expectedPreviouslyCheckedOutBranch, err := args.Runner.Backend.ExpectedPreviouslyCheckedOutBranch(step.InitialPreviouslyCheckedOutBranch, step.InitialBranch, step.MainBranch)
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
