package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// DeleteRemoteBranchStep deletes the tracking branch of the given local branch.
type DeleteRemoteBranchStep struct {
	Branch domain.RemoteBranchName
	EmptyStep
}

func (step *DeleteRemoteBranchStep) Run(args RunArgs) error {
	return args.Runner.Frontend.DeleteRemoteBranch(step.Branch)
}
