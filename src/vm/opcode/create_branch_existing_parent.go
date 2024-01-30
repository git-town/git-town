package opcode

import (
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/vm/shared"
)

// CreateBranchExistingParent creates a new branch with the first existing entry from the given ancestor list as its parent.
type CreateBranchExistingParent struct {
	Ancestors gitdomain.LocalBranchNames // list of ancestors - uses the first existing ancestor in this list
	Branch    gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *CreateBranchExistingParent) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		self,
	}
}

func (self *CreateBranchExistingParent) Run(args shared.RunArgs) error {
	nearestAncestor := args.Runner.Backend.FirstExistingBranch(self.Ancestors, args.Runner.MainBranch)
	return args.Runner.Frontend.CreateBranch(self.Branch, nearestAncestor.Location())
}
