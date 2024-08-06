package opcodes

import (
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/vm/shared"
)

// adds the branch with the given name as a prototype branch
type AddToPrototypeBranches struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *AddToPrototypeBranches) Run(args shared.RunArgs) error {
	return args.Config.AddToPrototypeBranches(self.Branch)
}
