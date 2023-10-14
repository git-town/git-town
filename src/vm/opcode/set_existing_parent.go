package opcode

import "github.com/git-town/git-town/v9/src/domain"

// SetExistingParent sets the first existing entry in the given ancestor list as the parent branch of the given branch.
type SetExistingParent struct {
	Branch     domain.LocalBranchName
	Ancestors  domain.LocalBranchNames
	MainBranch domain.LocalBranchName
	BaseOpcode
}

func (step *SetExistingParent) Run(args RunArgs) error {
	nearestAncestor := args.Runner.Backend.FirstExistingBranch(step.Ancestors, step.MainBranch)
	return args.Runner.Config.SetParent(step.Branch, nearestAncestor)
}
