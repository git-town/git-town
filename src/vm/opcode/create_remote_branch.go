package opcode

import (
	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/vm/shared"
)

// CreateRemoteBranch pushes the given local branch up to origin.
type CreateRemoteBranch struct {
	Branch     domain.LocalBranchName
	NoPushHook bool
	SHA        domain.SHA
	undeclaredOpcodeMethods
}

func (self *CreateRemoteBranch) Run(args shared.RunArgs) error {
	return args.Runner.Frontend.CreateRemoteBranch(self.SHA, self.Branch, self.NoPushHook)
}
