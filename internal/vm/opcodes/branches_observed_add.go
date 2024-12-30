package opcodes

import (
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/vm/shared"
)

// registers the branch with the given name as an observed branch in the Git config
// TODO: convert to more generic SetBranchTypeOverride opcode
type BranchesObservedAdd struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchesObservedAdd) Run(args shared.RunArgs) error {
	return args.Config.Value.NormalConfig.SetBranchTypeOverride(configdomain.BranchTypeObservedBranch, self.Branch)
}
