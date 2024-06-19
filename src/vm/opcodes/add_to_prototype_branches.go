package opcodes

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// AddToPerennialBranches adds the branch with the given name as a perennial branch.
type AddToPrototypeBranches struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *AddToPrototypeBranches) Run(args shared.RunArgs) error {
	return args.Config.AddToPrototypeBranches(self.Branch)
}
