package opcodes

import (
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/vm/shared"
)

// registers the branch with the given name as a parked branch in the Git config
// TODO: make generic
type BranchesParkedAdd struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchesParkedAdd) Run(args shared.RunArgs) error {
	return args.Config.Value.NormalConfig.SetBranchTypeOverride(configdomain.BranchTypeParkedBranch, self.Branch)
}
