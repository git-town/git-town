package steps

import (
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/runner"
)

// PreserveCheckoutHistoryStep does stuff
type PreserveCheckoutHistoryStep struct {
	NoOpStep
	InitialBranch                     string
	InitialPreviouslyCheckedOutBranch string
}

// Run executes this step.
func (step *PreserveCheckoutHistoryStep) Run() error {
	expectedPreviouslyCheckedOutBranch := git.GetExpectedPreviouslyCheckedOutBranch(step.InitialPreviouslyCheckedOutBranch, step.InitialBranch)
	if expectedPreviouslyCheckedOutBranch != git.GetPreviouslyCheckedOutBranch() {
		currentBranch := git.GetCurrentBranchName()
		runner.New("git", "checkout", expectedPreviouslyCheckedOutBranch).Run()
		runner.New("git", "checkout", currentBranch).Run()
	}
	return nil
}
