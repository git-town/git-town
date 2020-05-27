package steps

import (
	"github.com/git-town/git-town/src/git"
)

// PreserveCheckoutHistoryStep does stuff
type PreserveCheckoutHistoryStep struct {
	NoOpStep
	InitialBranch                     string
	InitialPreviouslyCheckedOutBranch string
}

// Run executes this step.
func (step *PreserveCheckoutHistoryStep) Run(repo *git.ProdRepo) error {
	expectedPreviouslyCheckedOutBranch, err := repo.Silent.GetExpectedPreviouslyCheckedOutBranch(step.InitialPreviouslyCheckedOutBranch, step.InitialBranch)
	if err != nil {
		return err
	}
	previousBranch, err := repo.Silent.PreviouslyCheckedOutBranch()
	if err != nil {
		return err
	}
	if expectedPreviouslyCheckedOutBranch != previousBranch {
		currentBranch, err := repo.Silent.CurrentBranch()
		if err != nil {
			return err
		}
		repo.Silent.CheckoutBranch(expectedPreviouslyCheckedOutBranch)
		repo.Silent.CheckoutBranch(currentBranch)
	}
	return nil
}
