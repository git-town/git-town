package opcode

import "github.com/git-town/git-town/v9/src/domain"

// CreateBranchExistingParent creates a new branch with the first existing entry from the given ancestor list as its parent.
type CreateBranchExistingParent struct {
	Branch     domain.LocalBranchName
	MainBranch domain.LocalBranchName
	Ancestors  domain.LocalBranchNames // list of ancestors - uses the first existing ancestor in this list
	BaseOpcode
}

func (step *CreateBranchExistingParent) Run(args RunArgs) error {
	nearestAncestor := args.Runner.Backend.FirstExistingBranch(step.Ancestors, step.MainBranch)
	return args.Runner.Frontend.CreateBranch(step.Branch, nearestAncestor.Location())
}
