package opcodes

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// SyncFeatureBranchMerge merges the parent branches of the given branch until a local parent is found.
type SyncFeatureBranchRebase struct {
	Branch               gitdomain.LocalBranchName
	ParentSHAPreviousRun Option[gitdomain.SHA]
	PushBranches         configdomain.PushBranches
	TrackingBranch       Option[gitdomain.RemoteBranchName]
}

func (self *SyncFeatureBranchRebase) Run(args shared.RunArgs) error {
	program := []shared.Opcode{}
	syncTracking, trackingBranchOpt, err := self.shouldSyncWithTracking(args)
	if err != nil {
		return err
	}
	trackingBranch, hasTrackingBranch := trackingBranchOpt.Get()
	if syncTracking && hasTrackingBranch {
		program = append(program,
			&RebaseTrackingBranch{
				PushBranches: self.PushBranches,
				RemoteBranch: trackingBranch,
			},
		)
	}
	program = append(program,
		&RebaseAncestorsUntilLocal{
			Branch:          self.Branch,
			CommitsToRemove: self.ParentSHAPreviousRun,
		},
	)
	// update the tracking branch
	if syncTracking && self.PushBranches.ShouldPush() && hasTrackingBranch && args.Config.Value.NormalConfig.Offline.IsOnline() {
		program = append(program,
			&PushCurrentBranchForceIfNeeded{
				CurrentBranch:   self.Branch,
				ForceIfIncludes: true,
			},
		)
	}
	args.PrependOpcodes(program...)
	return nil
}

func (self *SyncFeatureBranchRebase) shouldSyncWithTracking(args shared.RunArgs) (shouldSync bool, trackingBranchOpt Option[gitdomain.RemoteBranchName], err error) {
	trackingBranchOpt = self.TrackingBranch
	hasTrackingBranch := trackingBranchOpt.IsSome()
	if !hasTrackingBranch || args.Config.Value.NormalConfig.Offline.IsOffline() {
		return false, trackingBranchOpt, nil
	}
	syncedWithTracking, err := args.Git.BranchInSyncWithTracking(args.Backend, self.Branch, args.Config.Value.NormalConfig.DevRemote)
	return !syncedWithTracking, trackingBranchOpt, err
}
