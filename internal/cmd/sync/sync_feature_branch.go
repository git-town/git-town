package sync

import (
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/opcodes"
	"github.com/git-town/git-town/v16/internal/vm/program"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// FeatureBranchProgram adds the opcodes to sync the feature branch with the given name.
func FeatureBranchProgram(syncStrategy configdomain.SyncStrategy, args featureBranchArgs) {
	switch syncStrategy {
	case configdomain.SyncStrategyMerge:
		syncFeatureBranchMergeProgram(args)
	case configdomain.SyncStrategyRebase:
		syncFeatureBranchRebaseProgram(args)
	case configdomain.SyncStrategyCompress:
		syncFeatureBranchCompressProgram(args)
	}
}

type featureBranchArgs struct {
	firstCommitMessage Option[gitdomain.CommitMessage]
	localName          gitdomain.LocalBranchName
	offline            configdomain.Offline              // whether offline mode is enabled
	originalParentName Option[gitdomain.LocalBranchName] // the parent when Git Town started
	originalParentSHA  Option[gitdomain.SHA]             // the parent when Git Town started
	program            Mutable[program.Program]          // the program to update
	pushBranches       configdomain.PushBranches
	remoteName         Option[gitdomain.RemoteBranchName]
}

func syncFeatureBranchCompressProgram(args featureBranchArgs) {
	trackingBranch, hasTrackingBranch := args.remoteName.Get()
	if hasTrackingBranch {
		args.program.Value.Add(&opcodes.Merge{Branch: trackingBranch.BranchName()})
	}
	args.program.Value.Add(&opcodes.MergeParentIfNeeded{
		Branch:             args.localName,
		OriginalParentName: args.originalParentName,
		OriginalParentSHA:  args.originalParentSHA,
	})
	if firstCommitMessage, has := args.firstCommitMessage.Get(); has {
		args.program.Value.Add(&opcodes.BranchCurrentResetToParent{CurrentBranch: args.localName})
		args.program.Value.Add(&opcodes.CommitWithMessage{
			AuthorOverride: None[gitdomain.Author](),
			Message:        firstCommitMessage,
		})
	}
	if hasTrackingBranch && args.offline.IsFalse() {
		args.program.Value.Add(&opcodes.PushCurrentBranchForceIfNeeded{ForceIfIncludes: false})
	}
}

// syncs the given feature branch using the "merge" sync strategy
func syncFeatureBranchMergeProgram(args featureBranchArgs) {
	if trackingBranch, hasTrackingBranch := args.remoteName.Get(); hasTrackingBranch {
		args.program.Value.Add(&opcodes.Merge{Branch: trackingBranch.BranchName()})
	}
	args.program.Value.Add(&opcodes.MergeParentIfNeeded{
		Branch:             args.localName,
		OriginalParentName: args.originalParentName,
		OriginalParentSHA:  args.originalParentSHA,
	})
}

// syncs the given feature branch using the "rebase" sync strategy
func syncFeatureBranchRebaseProgram(args featureBranchArgs) {
	args.program.Value.Add(&opcodes.RebaseParentIfNeeded{
		Branch: args.localName,
	})
	if trackingBranch, hasTrackingBranch := args.remoteName.Get(); hasTrackingBranch {
		if args.offline.IsFalse() {
			args.program.Value.Add(&opcodes.RebaseTrackingBranch{RemoteBranch: trackingBranch, PushBranches: args.pushBranches})
		}
	}
}
