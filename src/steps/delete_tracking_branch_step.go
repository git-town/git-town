package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// DeleteTrackingBranchStep deletes the tracking branch of the given local branch.
type DeleteTrackingBranchStep struct {
	Branch domain.LocalBranchName
	EmptyStep
}

func (step *DeleteTrackingBranchStep) Run(args RunArgs) error {
	return args.Runner.Frontend.DeleteRemoteBranch(step.Branch)
}
