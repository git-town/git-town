package step

import "github.com/git-town/git-town/v9/src/domain"

// SetParent sets the given parent branch as the parent of the given branch.
// Use ChangeParent to change an existing parent.
type SetExistingParent struct {
	Branch     domain.LocalBranchName
	Ancestors  domain.LocalBranchNames
	MainBranch domain.LocalBranchName
	Empty
}

func (step *SetExistingParent) Run(args RunArgs) error {
	nearestAncestor := args.Runner.Backend.FirstExistingBranch(step.Ancestors, step.MainBranch)
	return args.Runner.Config.SetParent(step.Branch, nearestAncestor)
}
