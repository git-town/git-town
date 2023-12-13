package opcode

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/vm/shared"
)

// CreateRemoteBranch pushes the given local branch up to origin.
type CreateRemoteBranch struct {
	Branch     domain.LocalBranchName
	NoPushHook configdomain.NoPushHook
	SHA        domain.SHA
	undeclaredOpcodeMethods
}

func (self *CreateRemoteBranch) Run(args shared.RunArgs) error {
	return args.Runner.Frontend.CreateRemoteBranch(self.SHA, self.Branch, self.NoPushHook)
}
