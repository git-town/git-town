package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/runstate"
)

// IfBranchChangesStep executes different code paths
// depending on whether the given branch contains changes or not.
type IfBranchChangesStep struct {
	Branch          domain.LocalBranchName
	Parent          domain.LocalBranchName
	IsEmptySteps    runstate.StepList
	HasChangesSteps runstate.StepList
	EmptyStep
}

func (step *IfBranchChangesStep) Run(args RunArgs) error {
	hasChanges, err := args.Runner.Backend.BranchHasUnmergedChanges()
	if err != nil {
		return err
	}
	if hasChanges {

	}
}
