package step

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// IfBranchHasChanges executes different steps
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
		args.AddSteps(step.HasChangesSteps...)
	} else {
		args.AddSteps(step.IsEmptySteps...)
	}
	return nil
}
