package step

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// SetParent registers the branch with the given name as a parent
// of the branch with the other given name.
type SetParent struct {
	Branch       domain.LocalBranchName
	ParentBranch domain.LocalBranchName
	Empty
}

func (step *SetParent) Run(args RunArgs) error {
	return args.Runner.Config.SetParent(step.Branch, step.ParentBranch)
}
