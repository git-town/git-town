package opcodes

import (
	"github.com/git-town/git-town/v22/internal/config/gitconfig"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// BranchTypeOverrideRemove removes the branch with the given name from the contribution branches list in the Git config.
type BranchTypeOverrideRemove struct {
	Branch gitdomain.LocalBranchName
}

func (self *BranchTypeOverrideRemove) Run(args shared.RunArgs) error {
	_ = gitconfig.RemoveBranchTypeOverride(args.Backend, self.Branch)
	return nil
}
