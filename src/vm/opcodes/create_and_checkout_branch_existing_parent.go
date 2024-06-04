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
	currentBranch, err := args.Git.CurrentBranch(args.Backend)
	if err != nil {
		return err
	}
	nearestAncestor, hasNearestAncestor := args.Git.FirstExistingBranch(args.Backend, self.Ancestors...).Get()
	if !hasNearestAncestor {
		nearestAncestor = args.Config.Config.MainBranch
	}
	if nearestAncestor == currentBranch {
		return args.Git.CreateAndCheckoutBranch(args.Frontend, self.Branch)
	}
	return args.Git.CreateAndCheckoutBranchWithParent(args.Frontend, self.Branch, nearestAncestor.Location())
}
