package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// CreateRemoteBranch pushes the given local branch up to origin.
type CreateRemoteBranch struct {
	Branch     domain.LocalBranchName
	NoPushHook bool
	SHA        domain.SHA
	undeclaredOpcodeMethods
}

func (op *CreateRemoteBranch) Run(args shared.RunArgs) error {
	return args.Runner.Frontend.CreateRemoteBranch(op.SHA, op.Branch, op.NoPushHook)
}
