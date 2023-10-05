package step

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
)

// IfBranchHasChanges executes different code paths
// depending on whether the given branch contains changes or not.
type IfBranchHasChanges struct {
	Branch          domain.LocalBranchName
	Parent          domain.LocalBranchName
	IsEmptySteps    []Step
	HasChangesSteps []Step
	Empty
}

func (step *IfBranchHasChanges) Run(args RunArgs) error {
	hasChanges, err := args.Runner.Backend.BranchHasUnmergedChanges(step.Branch, step.Parent.Location())
	if err != nil {
		return err
	}
	if hasChanges {
		fmt.Println("111111111111111 BRANCH HAS CHANGES")
		args.AddSteps(step.HasChangesSteps...)
	} else {
		fmt.Println("111111111111111 BRANCH IS EMPTY")
		args.AddSteps(step.IsEmptySteps...)
	}
	return nil
}
