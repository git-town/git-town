package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// SetExistingParent sets the first existing entry in the given ancestor list as the parent branch of the given branch.
type SetExistingParent struct {
	Branch     domain.LocalBranchName
	Ancestors  domain.LocalBranchNames
	MainBranch domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (step *SetExistingParent) Run(args shared.RunArgs) error {
	nearestAncestor := args.Runner.Backend.FirstExistingBranch(step.Ancestors, step.MainBranch)
	return args.Runner.Config.SetParent(step.Branch, nearestAncestor)
}
