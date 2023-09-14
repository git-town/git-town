package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
)

// CreateTrackingBranchStep pushes the given local branch up to origin
// and marks it as tracking the current branch.
type CreateTrackingBranchStep struct {
	Branch     domain.LocalBranchName
	NoPushHook bool
	EmptyStep
}

func (step *CreateTrackingBranchStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&DeleteTrackingBranchStep{Branch: step.Branch, NoPushHook: false}}, nil
}

func (step *CreateTrackingBranchStep) Run(args RunArgs) error {
	return args.Run.Frontend.CreateTrackingBranch(step.Branch, domain.OriginRemote, step.NoPushHook)
}
