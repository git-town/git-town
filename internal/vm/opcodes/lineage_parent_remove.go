package opcodes

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

type LineageParentRemove struct {
	Branch gitdomain.LocalBranchName
}

func (self *LineageParentRemove) Run(args shared.RunArgs) error {
	if !args.Config.Value.NormalConfig.DryRun {
		args.Config.Value.NormalConfig.RemoveParent(args.Backend, self.Branch)
	}
	return nil
}
