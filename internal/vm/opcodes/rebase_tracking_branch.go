package opcodes

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// RebaseTrackingBranch rebases the current feature branch against its tracking branch.
type RebaseTrackingBranch struct {
	PushBranches configdomain.PushBranches
	RemoteBranch gitdomain.RemoteBranchName
}

func (self *RebaseTrackingBranch) Run(args shared.RunArgs) error {
	if self.PushBranches {
		// force-push-with-lease-if-includes here first, this avoids phantom merge conflicts from amended local commits
		err := args.Git.ForcePushBranchSafely(args.Frontend, args.Config.Value.NormalConfig.PushHook, true)
		if err == nil {
			// The force-push succeeded --> the remote branch didn't contain new commits, we are done.
			return nil
		}
	}
	args.PrependOpcodes(
		&RebaseBranch{Branch: self.RemoteBranch.BranchName()},
	)
	return nil
}
