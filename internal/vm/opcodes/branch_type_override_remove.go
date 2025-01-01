package opcodes

import (
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/vm/shared"
)

// removes the branch with the given name from the contribution branches list in the Git config
type BranchTypeOverrideRemove struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchTypeOverrideRemove) Run(args shared.RunArgs) error {
	return args.Config.Value.NormalConfig.RemoveBranchTypeOverride(self.Branch)
}
