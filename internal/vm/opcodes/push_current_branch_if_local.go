package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// PushCurrentBranch pushes the current branch to its existing tracking branch.
type PushCurrentBranchIfLocal struct {
	CurrentBranch           gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *PushCurrentBranchIfLocal) Run(args shared.RunArgs) error {
	branchIsLocal := args.Git.CurrentBranchHasTrackingBranch(args.Backend)
	if !branchIsLocal {
		args.PrependOpcodes(&PushCurrentBranch{
			CurrentBranch: self.CurrentBranch,
		})
	}
	return args.Git.PushCurrentBranch(args.Frontend, args.Config.Config.NoPushHook())
}
