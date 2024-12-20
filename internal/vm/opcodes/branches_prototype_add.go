package opcodes

import (
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/vm/shared"
)

// registers the branch with the given name as a prototype branch in the Git config
type BranchesPrototypeAdd struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchesPrototypeAdd) Run(args shared.RunArgs) error {
	return args.Config.Value.NormalConfig.AddToPrototypeBranches(self.Branch)
}
