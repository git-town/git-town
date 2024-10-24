package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// registers the branch with the given name as a perennial branch in the Git config
type BranchesPerennialAdd struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchesPerennialAdd) Run(args shared.RunArgs) error {
	return args.Config.Value.NormalConfig.AddToPerennialBranches(self.Branch)
}
