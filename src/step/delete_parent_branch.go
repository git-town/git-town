package step

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// DeleteParentBranch removes the parent branch entry in the Git Town configuration.
type DeleteParentBranch struct {
	Branch domain.LocalBranchName
	Empty
}

func (step *DeleteParentBranch) Run(args RunArgs) error {
	return args.Runner.Config.RemoveParent(step.Branch)
}
