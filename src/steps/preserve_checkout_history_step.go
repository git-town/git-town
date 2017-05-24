package steps

import (
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/util"
)

// PreserveCheckoutHistoryStep does stuff
type PreserveCheckoutHistoryStep struct {
	NoOpStep
	InitialBranch                     string
	InitialPreviouslyCheckedOutBranch string
}

// Run executes this step.
func (step PreserveCheckoutHistoryStep) Run() error {
	expectedPreviouslyCheckedOutBranch := git.GetExpectedPreviouslyCheckedOutBranch(step.InitialPreviouslyCheckedOutBranch, step.InitialBranch)
	if expectedPreviouslyCheckedOutBranch != git.GetPreviouslyCheckedOutBranch() {
		currentBranch := git.GetCurrentBranchName()
		util.GetCommandOutput("git", "checkout", expectedPreviouslyCheckedOutBranch)
		util.GetCommandOutput("git", "checkout", currentBranch)
	}
	return nil
}
