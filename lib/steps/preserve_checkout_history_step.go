package steps

import (
	"fmt"

	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/util"
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
	fmt.Println("initial:", step.InitialPreviouslyCheckedOutBranch)
	fmt.Println("expected:", expectedPreviouslyCheckedOutBranch)
	fmt.Println("current:", git.GetPreviouslyCheckedOutBranch())
	if expectedPreviouslyCheckedOutBranch != git.GetPreviouslyCheckedOutBranch() {
		currentBranch := git.GetCurrentBranchName()
		util.GetCommandOutput("git", "checkout", expectedPreviouslyCheckedOutBranch)
		util.GetCommandOutput("git", "checkout", currentBranch)
	}
	return nil
}
