package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

type LineageParentRemove struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *LineageParentRemove) Run(args shared.RunArgs) error {
	args.Config.RemoveParent(self.Branch)
	return nil
}
