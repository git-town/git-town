package opcodes

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

// PushCurrentBranchForceIfNeeded force-pushes the branch with the given name to the origin remote.
type PushCurrentBranchForceIfNeeded struct {
	CurrentBranch           gitdomain.LocalBranchName
	ForceIfIncludes         bool
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *PushCurrentBranchForceIfNeeded) Run(args shared.RunArgs) error {
	inSync, err := args.Git.BranchInSyncWithTracking(args.Backend, self.CurrentBranch, args.Config.Value.NormalConfig.DevRemote)
	if err != nil {
		return err
	}
	if inSync {
		return nil
	}
	args.PrependOpcodes(&PushCurrentBranchForce{ForceIfIncludes: self.ForceIfIncludes})
	return nil
}
