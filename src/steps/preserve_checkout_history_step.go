package steps

import (
	"github.com/Originate/git-town/src/command"
	"github.com/Originate/git-town/src/git"
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
		command.MustRun("git", "checkout", expectedPreviouslyCheckedOutBranch)
		command.MustRun("git", "checkout", currentBranch)
	}
	return nil
}
