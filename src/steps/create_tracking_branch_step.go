package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// CreateTrackingBranchStep pushes the given local branch up to origin
// and marks it as tracking the current branch.
type CreateTrackingBranchStep struct {
	Branch     domain.LocalBranchName
	NoPushHook bool
	EmptyStep
}

func (step *CreateTrackingBranchStep) Run(args RunArgs) error {
	return args.Runner.Frontend.CreateTrackingBranch(step.Branch, domain.OriginRemote, step.NoPushHook)
}
