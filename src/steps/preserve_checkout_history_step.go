package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// PreserveCheckoutHistoryStep does stuff.
type PreserveCheckoutHistoryStep struct {
	InitialBranch                     domain.LocalBranchName
	InitialPreviouslyCheckedOutBranch domain.LocalBranchName
	MainBranch                        domain.LocalBranchName
	EmptyStep
}

func (step *PreserveCheckoutHistoryStep) Run(args RunArgs) error {
	expectedPreviouslyCheckedOutBranch, err := args.Run.Backend.ExpectedPreviouslyCheckedOutBranch(step.InitialPreviouslyCheckedOutBranch, step.InitialBranch, step.MainBranch)
	if err != nil {
		return err
	}
	if expectedPreviouslyCheckedOutBranch == args.Run.Backend.PreviouslyCheckedOutBranch() {
		return nil
	}
	currentBranch, err := args.Run.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	err = args.Run.Backend.CheckoutBranchUncached(expectedPreviouslyCheckedOutBranch)
	if err != nil {
		return err
	}
	return args.Run.Backend.CheckoutBranchUncached(currentBranch)
}
