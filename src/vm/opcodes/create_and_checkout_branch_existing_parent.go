package opcodes

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// CreateAndCheckoutBranchExistingParent creates a new branch with the first existing entry from the given ancestor list as its parent.
type CreateAndCheckoutBranchExistingParent struct {
	Ancestors               gitdomain.LocalBranchNames // list of ancestors - uses the first existing ancestor in this list
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CreateAndCheckoutBranchExistingParent) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		self,
	}
}

func (self *CreateAndCheckoutBranchExistingParent) Run(args shared.RunArgs) error {
	nearestAncestor := args.Backend.FirstExistingBranch(self.Ancestors, args.Config.Config.MainBranch)
	currentBranch, err := args.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	if nearestAncestor == currentBranch {
		return args.Frontend.CreateAndCheckoutBranch(self.Branch)
	}
	return args.Frontend.CreateAndCheckoutBranchWithParent(self.Branch, nearestAncestor.Location())
}
