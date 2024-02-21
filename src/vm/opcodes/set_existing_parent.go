package opcodes

import (
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/vm/shared"
)

// SetExistingParent sets the first existing entry in the given ancestor list as the parent branch of the given branch.
type SetExistingParent struct {
	Ancestors gitdomain.LocalBranchNames
	Branch    gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *SetExistingParent) Run(args shared.RunArgs) error {
	nearestAncestor := args.Runner.Backend.FirstExistingBranch(self.Ancestors, args.Runner.Config.FullConfig.MainBranch)
	return args.Runner.Config.SetParent(self.Branch, nearestAncestor)
}
