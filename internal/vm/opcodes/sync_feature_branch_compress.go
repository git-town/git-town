package opcodes

import (
	"github.com/git-town/git-town/v20/internal/config/configdomain"
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/vm/shared"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

// SyncFeatureBranchCompress expands to all opcodes needed to sync a feature branch using the "compress" sync strategy.
type SyncFeatureBranchCompress struct {
	CommitMessage           Option[gitdomain.CommitMessage]
	CurrentBranch           gitdomain.LocalBranchName
	InitialParentName       Option[gitdomain.LocalBranchName]
	InitialParentSHA        Option[gitdomain.SHA]
	Offline                 configdomain.Offline
	TrackingBranch          Option[gitdomain.RemoteBranchName]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *SyncFeatureBranchCompress) Run(args shared.RunArgs) error {
	opcodes := []shared.Opcode{
		&MergeParentsUntilLocal{
			Branch:            self.CurrentBranch,
			InitialParentName: self.InitialParentName,
			InitialParentSHA:  self.InitialParentSHA,
		},
	}
	if trackingBranch, hasTrackingBranch := self.TrackingBranch.Get(); hasTrackingBranch {
		opcodes = append(opcodes,
			&MergeIntoCurrentBranch{BranchToMerge: trackingBranch.BranchName()},
		)
		if firstCommitMessage, has := self.CommitMessage.Get(); has {
			opcodes = append(opcodes, &BranchCurrentResetToParent{CurrentBranch: self.CurrentBranch})
			opcodes = append(opcodes, &CommitWithMessage{
				AuthorOverride: None[gitdomain.Author](),
				CommitHook:     configdomain.CommitHookEnabled,
				Message:        firstCommitMessage,
			})
		}
		if self.Offline.IsFalse() {
			opcodes = append(opcodes, &PushCurrentBranchForceIfNeeded{CurrentBranch: self.CurrentBranch, ForceIfIncludes: false})
		}
	}
	args.PrependOpcodes(opcodes...)
	return nil
}
