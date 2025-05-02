package opcodes

import (
	"github.com/git-town/git-town/v20/internal/config/configdomain"
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/vm/shared"
)

// RebaseTrackingBranch rebases the current feature branch against its tracking branch.
type RebaseTrackingBranch struct {
	PushBranches            configdomain.PushBranches
	RemoteBranch            gitdomain.RemoteBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseTrackingBranch) Run(args shared.RunArgs) error {
	// Try to force-push the local branch with lease and includes to the remote branch.
	if self.PushBranches {
		err := args.Git.ForcePushBranchSafely(args.Frontend, args.Config.Value.NormalConfig.NoPushHook(), true)
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
