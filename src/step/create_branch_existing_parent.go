package step

import "github.com/git-town/git-town/v9/src/domain"

// CreateBranchExistingParent cuts a new branch from the given starting point,
// or if it doesn't exist its nearest existing ancestor.
type CreateBranchExistingParent struct {
	Branch        domain.LocalBranchName
	MainBranch    domain.LocalBranchName
	StartingPoint domain.LocalBranchName
	Ancestors     domain.LocalBranchNames // list of ancestors - uses the first existing ancestor in this list
	Empty
}

func (step *CreateBranchExistingParent) Run(args RunArgs) error {
	nearestAncestor := args.Runner.Backend.FirstExistingBranch(step.Ancestors, step.MainBranch)
	return args.Runner.Frontend.CreateBranch(step.Branch, nearestAncestor.Location())
}
