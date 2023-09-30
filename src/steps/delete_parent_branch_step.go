package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// DeleteParentBranchStep removes the parent branch entry in the Git Town configuration.
type DeleteParentBranchStep struct {
	Branch domain.LocalBranchName
	EmptyStep
}

func (step *DeleteParentBranchStep) Run(args RunArgs) error {
	return args.Runner.Config.RemoveParent(step.Branch)
}
