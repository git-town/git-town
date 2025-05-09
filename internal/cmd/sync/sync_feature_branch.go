package sync

import (
	"github.com/git-town/git-town/v20/internal/config/configdomain"
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/vm/opcodes"
	"github.com/git-town/git-town/v20/internal/vm/program"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

// FeatureBranchProgram adds the opcodes to sync the feature branch with the given name.
func FeatureBranchProgram(syncStrategy configdomain.SyncStrategy, args featureBranchArgs) {
	switch syncStrategy {
	case configdomain.SyncStrategyCompress:
		syncFeatureBranchCompress(args)
	case configdomain.SyncStrategyFFOnly:
		syncFeatureBranchFFOnly(args)
	case configdomain.SyncStrategyMerge:
		syncFeatureBranchMerge(args)
	case configdomain.SyncStrategyRebase:
		syncFeatureBranchRebase(args)
	}
	if args.prune {
		args.program.Value.Add(&opcodes.BranchDeleteIfEmptyAtRuntime{Branch: args.localName})
	}
}

type featureBranchArgs struct {
	firstCommitMessage Option[gitdomain.CommitMessage]
	localName          gitdomain.LocalBranchName
	offline            configdomain.Offline              // whether offline mode is enabled
	originalParentName Option[gitdomain.LocalBranchName] // the parent when Git Town started
	originalParentSHA  Option[gitdomain.SHA]             // the parent when Git Town started
	parentLastRunSHA   Option[gitdomain.SHA]             // the parent at the end of the last Git Town command
	program            Mutable[program.Program]          // the program to update
	prune              configdomain.Prune
	pushBranches       configdomain.PushBranches
	trackingBranchName Option[gitdomain.RemoteBranchName]
}

func syncFeatureBranchCompress(args featureBranchArgs) {
	// sync parent branch
	args.program.Value.Add(
		&opcodes.MergeParentsUntilLocal{
			Branch:             args.localName,
			OriginalParentName: args.originalParentName,
			OriginalParentSHA:  args.originalParentSHA,
		},
	)
	// sync tracking branch
	if trackingBranch, hasTrackingBranch := args.trackingBranchName.Get(); hasTrackingBranch {
		args.program.Value.Add(
			&opcodes.CompressMergeTrackingBranch{
				CurrentBranch:  args.localName,
				CommitMessage:  args.firstCommitMessage,
				Offline:        args.offline,
				TrackingBranch: trackingBranch,
			},
		)
	}
}

func syncFeatureBranchFFOnly(args featureBranchArgs) {
	// The ff-only strategy does not sync with the parent branch.
	// It is intended for perennial branches only.
	if args.offline.IsFalse() {
		if trackingBranch, hasTrackingBranch := args.trackingBranchName.Get(); hasTrackingBranch {
			args.program.Value.Add(&opcodes.MergeFastForward{Branch: trackingBranch.BranchName()})
		}
	}
}

func syncFeatureBranchMerge(args featureBranchArgs) {
	// sync parent branch
	args.program.Value.Add(
		&opcodes.MergeParentsUntilLocal{
			Branch:             args.localName,
			OriginalParentName: args.originalParentName,
			OriginalParentSHA:  args.originalParentSHA,
		},
	)
	// sync tracking branch
	if trackingBranch, hasTrackingBranch := args.trackingBranchName.Get(); hasTrackingBranch {
		args.program.Value.Add(&opcodes.MergeIntoCurrentBranch{BranchToMerge: trackingBranch.BranchName()})
	}
}

func syncFeatureBranchRebase(args featureBranchArgs) {
	// sync parent branch
	args.program.Value.Add(
		&opcodes.RebaseParentsUntilLocal{
			Branch:      args.localName,
			PreviousSHA: args.parentLastRunSHA,
		},
	)
	// sync tracking branch
	if trackingBranch, hasTrackingBranch := args.trackingBranchName.Get(); hasTrackingBranch {
		if args.offline.IsFalse() {
			args.program.Value.Add(
				&opcodes.RebaseTrackingBranch{
					RemoteBranch: trackingBranch,
					PushBranches: args.pushBranches,
				},
				&opcodes.RebaseParentsUntilLocal{
					Branch:      args.localName,
					PreviousSHA: args.parentLastRunSHA,
				},
				&opcodes.PushCurrentBranchForceIfNeeded{
					CurrentBranch:   args.localName,
					ForceIfIncludes: true,
				},
			)
		}
	}
}
