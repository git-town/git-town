package opcodes

import (
	"github.com/git-town/git-town/v20/internal/config/configdomain"
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/vm/shared"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

// SyncFeatureBranchMerge merges the parent branches of the given branch until a local parent is found.
type SyncFeatureBranchRebase struct {
	Branch gitdomain.LocalBranchName
	// InitialParentName       Option[gitdomain.LocalBranchName]
	// InitialParentSHA        Option[gitdomain.SHA]
	ParentLastRunSHA        Option[gitdomain.SHA]
	PushBranches            configdomain.PushBranches
	TrackingBranch          Option[gitdomain.RemoteBranchName]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *SyncFeatureBranchRebase) Run(args shared.RunArgs) error {
	// new sync workflow:
	// 1. determine which branches need to be synced
	//    - parent branch
	//    - tracking branch
	// 2. if either parent or tracking branch needs to be synced: sync both, then force-push to tracking

	syncTracking := false
	trackingBranch, hasTrackingBranch := self.TrackingBranch.Get()
	if hasTrackingBranch {
		if args.Config.Value.NormalConfig.Offline.IsOnline() {
			syncedWithTracking, err := args.Git.BranchInSyncWithTracking(args.Backend, self.Branch, args.Config.Value.NormalConfig.DevRemote)
			if err != nil {
				return err
			}
			syncTracking = !syncedWithTracking
		}
	}
	program := []shared.Opcode{}
	if syncTracking {
		program = append(program,
			&RebaseTrackingBranch{
				PushBranches: self.PushBranches,
				RemoteBranch: trackingBranch,
			},
		)
	}
	program = append(program,
		&RebaseParentsUntilLocal{
			Branch:      self.Branch,
			PreviousSHA: self.ParentLastRunSHA,
		},
	)
	// update the tracking branch
	if self.PushBranches.IsTrue() && hasTrackingBranch {
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
