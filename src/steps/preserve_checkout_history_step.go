package steps

import (
	"github.com/git-town/git-town/v7/src/drivers"
	"github.com/git-town/git-town/v7/src/git"
)

// PreserveCheckoutHistoryStep does stuff.
type PreserveCheckoutHistoryStep struct {
	NoOpStep
	InitialBranch                     string
	InitialPreviouslyCheckedOutBranch string
}

// Run executes this step.
func (step *PreserveCheckoutHistoryStep) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	expectedPreviouslyCheckedOutBranch, err := repo.Silent.ExpectedPreviouslyCheckedOutBranch(step.InitialPreviouslyCheckedOutBranch, step.InitialBranch)
	if err != nil {
		return err
	}
	// NOTE: errors are not a failure condition here --> ignoring them
	previouslyCheckedOutBranch, _ := repo.Silent.PreviouslyCheckedOutBranch()
	if expectedPreviouslyCheckedOutBranch == previouslyCheckedOutBranch {
		return nil
	}
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
