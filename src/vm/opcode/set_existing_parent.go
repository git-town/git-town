package opcode

import (
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/vm/shared"
)

// SetExistingParent sets the first existing entry in the given ancestor list as the parent branch of the given branch.
type SetExistingParent struct {
	Branch     gitdomain.LocalBranchName
	Ancestors  gitdomain.LocalBranchNames
	MainBranch gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *SetExistingParent) Run(args shared.RunArgs) error {
	nearestAncestor := args.Runner.Backend.FirstExistingBranch(self.Ancestors, self.MainBranch)
	return args.Runner.GitTown.SetParent(self.Branch, nearestAncestor)
}
