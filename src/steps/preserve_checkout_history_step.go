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
	expectedPreviouslyCheckedOutBranch, err := repo.Silent.ExpectedPreviouslyCheckedOutBranch(step.InitialPreviouslyCheckedOutBranch, step.InitialBranch)
	if err != nil {
		return err
	}
	if expectedPreviouslyCheckedOutBranch != git.GetPreviouslyCheckedOutBranch() {
		currentBranch, err := repo.Silent.CurrentBranch()
		if err != nil {
			return err
		}
		err = repo.Silent.CheckoutBranch(expectedPreviouslyCheckedOutBranch)
		if err != nil {
			return err
		}
		return repo.Silent.CheckoutBranch(currentBranch)
	}
	return nil
}
