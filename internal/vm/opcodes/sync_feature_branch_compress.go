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
	Offline                 configdomain.Offline
	ParentName              Option[gitdomain.LocalBranchName]
	ParentSHA               Option[gitdomain.SHA]
	TrackingBranch          Option[gitdomain.RemoteBranchName]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *SyncFeatureBranchCompress) Run(args shared.RunArgs) error {
	opcodes := []shared.Opcode{}
	if parentName, hasParent := self.ParentName.Get(); hasParent {
		inSyncWithParent, err := args.Git.BranchInSyncWithParent(args.Backend, self.CurrentBranch, parentName)
		if err != nil {
			return err
		}
		if !inSyncWithParent {
			opcodes = append(opcodes, &MergeParentsUntilLocal{
				Branch:             self.CurrentBranch,
				OriginalParentName: self.ParentName,
				OriginalParentSHA:  self.ParentSHA,
			})
		}
	}
	if trackingBranch, hasTrackingBranch := self.TrackingBranch.Get(); hasTrackingBranch {
		opcodes = append(opcodes,
			&MergeIntoCurrentBranch{BranchToMerge: trackingBranch.BranchName()},
		)
	}
	if len(opcodes) > 0 {
		firstCommitMessage := self.CommitMessage.GetOrElse("compressed commit")
		opcodes = append(opcodes,
			&BranchCurrentResetToParent{
				CurrentBranch: self.CurrentBranch,
			},
			&CommitWithMessage{
				AuthorOverride: None[gitdomain.Author](),
				CommitHook:     configdomain.CommitHookEnabled,
				Message:        firstCommitMessage,
			},
		)
		if self.Offline.IsFalse() {
			opcodes = append(opcodes, &PushCurrentBranchForceIfNeeded{CurrentBranch: self.CurrentBranch, ForceIfIncludes: false})
		}
	}
	args.PrependOpcodes(opcodes...)
	return nil
}
