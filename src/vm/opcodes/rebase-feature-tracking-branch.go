package opcodes

import (
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/vm/shared"
)

// PushFeatureTrackingBranch pushes newly created Git tags to origin.
type RebaseFeatureTrackingBranch struct {
	remoteBranch gitdomain.RemoteBranchName
	undeclaredOpcodeMethods
}

func (self *RebaseFeatureTrackingBranch) Run(args shared.RunArgs) error {
	err := args.Runner.Frontend.ForcePushBranch(args.Runner.Config.FullConfig.NoPushHook())
	if err == nil {
		return nil
	}
	args.PrependOpcodes(
		&RebaseBranch{Branch: self.remoteBranch.BranchName()},
		&RebaseFeatureTrackingBranch{},
	)
	return nil
}
