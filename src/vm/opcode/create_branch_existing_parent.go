package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// CreateBranchExistingParent creates a new branch with the first existing entry from the given ancestor list as its parent.
type CreateBranchExistingParent struct {
	Branch     domain.LocalBranchName
	MainBranch domain.LocalBranchName
	Ancestors  domain.LocalBranchNames // list of ancestors - uses the first existing ancestor in this list
	undeclaredOpcodeMethods
}

func (op *CreateBranchExistingParent) Run(args shared.RunArgs) error {
	nearestAncestor := args.Runner.Backend.FirstExistingBranch(op.Ancestors, op.MainBranch)
	return args.Runner.Frontend.CreateBranch(op.Branch, nearestAncestor.Location())
}
