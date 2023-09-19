package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
)

// DeleteTrackingBranchStep deletes the tracking branch of the given local branch.
type DeleteTrackingBranchStep struct {
	Branch     domain.LocalBranchName
	NoPushHook bool
	EmptyStep
}

func (step *DeleteTrackingBranchStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&CreateTrackingBranchStep{Branch: step.Branch, NoPushHook: step.NoPushHook}}, nil
}

func (step *DeleteTrackingBranchStep) Run(args RunArgs) error {
	return args.Runner.Frontend.DeleteRemoteBranch(step.Branch)
}
