package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// PreserveCheckoutHistoryStep does stuff.
type PreserveCheckoutHistoryStep struct {
	EmptyStep
	InitialBranch                     string
	InitialPreviouslyCheckedOutBranch string
}

func (step *PreserveCheckoutHistoryStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
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
