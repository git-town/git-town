package opcodes

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// PushCurrentBranchForceIfNeeded force-pushes the branch with the given name to the origin remote.
type PushCurrentBranchForceIfNeeded struct {
	CurrentBranch   gitdomain.LocalBranchName
	TrackingBranch  gitdomain.RemoteBranchName
	ForceIfIncludes bool
}

func (self *PushCurrentBranchForceIfNeeded) Run(args shared.RunArgs) error {
	inSync, err := args.Git.BranchInSyncWithTracking(args.Backend, self.CurrentBranch, self.TrackingBranch)
	if err != nil {
		return err
	}
	if inSync {
		return nil
	}
	args.PrependOpcodes(&PushCurrentBranchForce{ForceIfIncludes: self.ForceIfIncludes})
	return nil
}
