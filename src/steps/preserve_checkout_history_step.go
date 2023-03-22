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
	MainBranch                        string
}

func (step *PreserveCheckoutHistoryStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	expectedPreviouslyCheckedOutBranch, err := repo.Internal.ExpectedPreviouslyCheckedOutBranch(step.InitialPreviouslyCheckedOutBranch, step.InitialBranch, step.MainBranch)
	if err != nil {
		return err
	}
	// NOTE: errors are not a failure condition here --> ignoring them
	previouslyCheckedOutBranch, _ := repo.Internal.PreviouslyCheckedOutBranch()
	if expectedPreviouslyCheckedOutBranch == previouslyCheckedOutBranch {
		return nil
	}
	currentBranch, err := repo.Internal.CurrentBranch()
	if err != nil {
		return err
	}
	err = repo.Internal.CheckoutBranch(expectedPreviouslyCheckedOutBranch)
	if err != nil {
		return err
	}
	return repo.Internal.CheckoutBranch(currentBranch)
}
