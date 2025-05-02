package opcodes

import (
	"github.com/git-town/git-town/v20/internal/config/configdomain"
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/vm/shared"
)

// registers the branch with the given name as a contribution branch in the Git config
type BranchTypeOverrideSet struct {
	Branch                  gitdomain.LocalBranchName
	BranchType              configdomain.BranchType
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchTypeOverrideSet) Run(args shared.RunArgs) error {
	return args.Config.Value.NormalConfig.SetBranchTypeOverride(self.BranchType, self.Branch)
}
