package opcodes

import (
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/vm/shared"
)

// RebaseFeatureTrackingBranch rebases the current feature branch against its tracking branch.
type RebaseFeatureTrackingBranch struct {
	PushBranches            configdomain.PushBranches
	RemoteBranch            gitdomain.RemoteBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseFeatureTrackingBranch) Run(args shared.RunArgs) error {
	// Try to force-push the local branch with lease and includes to the remote branch.
	if self.PushBranches {
		err := args.Git.ForcePushBranchSafely(args.Frontend, args.Config.Config.NoPushHook(), true)
		if err == nil {
			// The force-push succeeded --> the remote branch didn't contain new commits, we are done.
			return nil
		}
		// Here the force-push failed --> the remote branch contains new commits.
		// We need to integrate them into the local branch.
		args.PrependOpcodes(
			// Rebase the local commits against the remote commits.
			&RebaseBranch{Branch: self.RemoteBranch.BranchName()},
			// Now try force-pushing again.
			&RebaseFeatureTrackingBranch{
				PushBranches: self.PushBranches,
				RemoteBranch: self.RemoteBranch,
			},
		)
	} else {
		args.PrependOpcodes(&RebaseBranch{Branch: self.RemoteBranch.BranchName()})
	}
	return nil
}
