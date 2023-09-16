package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// SetParentStep registers the branch with the given name as a parent
// of the branch with the other given name.
type SetParentStep struct {
	Branch       domain.LocalBranchName
	ParentBranch domain.LocalBranchName
	EmptyStep
}

func (step *SetParentStep) Run(args RunArgs) error {
	return args.Runner.Config.SetParent(step.Branch, step.ParentBranch)
}
