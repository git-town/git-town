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
	DevRemote               gitdomain.Remote
	Offline                 configdomain.Offline
	TrackingBranch          gitdomain.RemoteBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CompressMergeTrackingBranch) Run(args shared.RunArgs) error {
	isInSync, err := args.Git.BranchInSyncWithTracking(args.Backend, self.CurrentBranch, self.DevRemote)
	if err != nil || isInSync {
		return err
	}
	opcodes := []shared.Opcode{
		&MergeIntoCurrentBranch{BranchToMerge: self.TrackingBranch.BranchName()},
	}
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
	args.PrependOpcodes(opcodes...)
	return nil
}
