package opcodes

import (
	"github.com/git-town/git-town/v20/internal/config/configdomain"
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/vm/shared"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

// CompressMergeTrackingBranch merges the tracking branch when syncing using the "compress" strategy.
type CompressMergeTrackingBranch struct {
	CommitMessage           Option[gitdomain.CommitMessage]
	CurrentBranch           gitdomain.LocalBranchName
	Offline                 configdomain.Offline
	ParentName              Option[gitdomain.LocalBranchName]
	ParentSHA               Option[gitdomain.SHA]
	TrackingBranch          Option[gitdomain.RemoteBranchName]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CompressMergeTrackingBranch) Run(args shared.RunArgs) error {
	opcodes := []shared.Opcode{
		&MergeParentsUntilLocal{
			Branch:             self.CurrentBranch,
			OriginalParentName: self.ParentName,
			OriginalParentSHA:  self.ParentSHA,
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
