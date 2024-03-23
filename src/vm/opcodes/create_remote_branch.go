package opcodes

import (
	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/git-town/git-town/v13/src/vm/shared"
)

// CreateRemoteBranch pushes the given local branch up to origin.
type CreateRemoteBranch struct {
	Branch gitdomain.LocalBranchName
	SHA    gitdomain.SHA
	undeclaredOpcodeMethods
}

func (self *CreateRemoteBranch) Run(args shared.RunArgs) error {
	return args.Runner.Frontend.CreateRemoteBranch(self.SHA, self.Branch, args.Runner.Config.FullConfig.NoPushHook())
}
