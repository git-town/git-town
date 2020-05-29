package steps

import (
	"github.com/git-town/git-town/src/command"
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
		currentBranch := git.GetCurrentBranchName()
		command.MustRun("git", "checkout", expectedPreviouslyCheckedOutBranch)
		command.MustRun("git", "checkout", currentBranch)
	}
	return nil
}
