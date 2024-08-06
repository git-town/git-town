package opcodes

import (
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/vm/shared"
)

// CreateRemoteBranch pushes the given local branch up to origin.
type CreateRemoteBranch struct {
	Branch                  gitdomain.LocalBranchName
	SHA                     gitdomain.SHA
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CreateRemoteBranch) Run(args shared.RunArgs) error {
	return args.Git.CreateRemoteBranch(args.Frontend, self.SHA, self.Branch, args.Config.Config.NoPushHook())
}
