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
	program := []shared.Opcode{
		&RebaseParentsUntilLocal{
			Branch:      self.Branch,
			PreviousSHA: self.ParentLastRunSHA,
		},
	}
	if trackingBranch, hasTrackingBranch := self.TrackingBranch.Get(); hasTrackingBranch {
		if args.Config.Value.NormalConfig.Offline.IsOnline() {
			program = append(program,
				&RebaseTrackingBranch{
					PushBranches: self.PushBranches,
					RemoteBranch: trackingBranch,
				},
				&RebaseParentsUntilLocal{
					Branch:      self.Branch,
					PreviousSHA: self.ParentLastRunSHA,
				},
				&PushCurrentBranchForceIfNeeded{
					CurrentBranch:   self.Branch,
					ForceIfIncludes: true,
				},
			)
		}
	}
	args.PrependOpcodes(program...)
	return nil
}
