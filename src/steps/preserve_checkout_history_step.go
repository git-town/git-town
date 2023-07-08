package steps

import (
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// PreserveCheckoutHistoryStep does stuff.
type PreserveCheckoutHistoryStep struct {
	EmptyStep
	InitialBranch                     string
	InitialPreviouslyCheckedOutBranch string
	MainBranch                        string
}

func (step *PreserveCheckoutHistoryStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	expectedPreviouslyCheckedOutBranch, err := run.Backend.ExpectedPreviouslyCheckedOutBranch(step.InitialPreviouslyCheckedOutBranch, step.InitialBranch, step.MainBranch)
	if err != nil {
		return err
	}
	// NOTE: errors are not a failure condition here --> ignoring them
	previouslyCheckedOutBranch, _ := run.Backend.PreviouslyCheckedOutBranch()
	if expectedPreviouslyCheckedOutBranch == previouslyCheckedOutBranch {
		return nil
	}
	currentBranch, err := run.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	err = run.Backend.CheckoutBranch(expectedPreviouslyCheckedOutBranch)
	if err != nil {
		return err
	}
	return run.Backend.CheckoutBranch(currentBranch)
}
