package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

type BranchParentDelete struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchParentDelete) Run(args shared.RunArgs) error {
	args.Config.Value.NormalConfig.RemoveParent(self.Branch)
	return nil
}
