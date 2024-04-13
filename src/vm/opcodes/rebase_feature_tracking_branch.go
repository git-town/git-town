package opcodes

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// RebaseFeatureTrackingBranch rebases the current feature branch against its tracking branch.
type RebaseFeatureTrackingBranch struct {
	RemoteBranch gitdomain.RemoteBranchName
	undeclaredOpcodeMethods
}

func (self *RebaseFeatureTrackingBranch) Run(args shared.RunArgs) error {
	// Try to force-push the local branch with lease and includes to the remote branch.
	err := args.Runner.Frontend.ForcePushBranchSafely(args.Runner.Config.FullConfig.NoPushHook())
	if err == nil {
		// The force-push succeeded --> the remote branch didn't contain new commits, we are done.
		return nil
	}
	// The force-push failed --> the remote branch contains new commits.
	// We need to integrate them into the local branch.
	args.PrependOpcodes(
		// Rebase the local commits against the remote commits.
		&RebaseBranch{Branch: self.RemoteBranch.BranchName()},
		// Now try force-pushing again.
		&RebaseFeatureTrackingBranch{RemoteBranch: self.RemoteBranch},
	)
	return nil
}
