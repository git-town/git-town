package opcodes

import (
	"github.com/git-town/git-town/v14/internal/git/gitdomain"
	"github.com/git-town/git-town/v14/internal/vm/shared"
)

// AddToPerennialBranches adds the branch with the given name as a perennial branch.
type AddToPerennialBranches struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *AddToPerennialBranches) Run(args shared.RunArgs) error {
	return args.Config.AddToPerennialBranches(self.Branch)
}
