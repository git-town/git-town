package opcodes

import (
	"github.com/git-town/git-town/v18/internal/config/configdomain"
	"github.com/git-town/git-town/v18/internal/git/gitdomain"
	"github.com/git-town/git-town/v18/internal/vm/shared"
)

// RebaseTrackingBranch rebases the current feature branch against its tracking branch.
type RebaseTrackingBranch struct {
	PushBranches            configdomain.PushBranches
	RemoteBranch            gitdomain.RemoteBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseTrackingBranch) Run(args shared.RunArgs) error {
	err := args.Git.Rebase(args.Frontend, self.RemoteBranch.BranchName(), args.Config.Value.NormalConfig.GitVersion)
	if err != nil {
		return err
	}
	if self.PushBranches {
		// ignoring push errors here - pushes can fail if the branch is in the merge queue
		_ = args.Git.ForcePushBranchSafely(args.Frontend, args.Config.Value.NormalConfig.NoPushHook(), true)
	}
	return nil
}
