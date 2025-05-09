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
	syncFeatureParentBranch(syncStrategy, args)
	if trackingBranch, hasTrackingBranch := args.trackingBranchName.Get(); hasTrackingBranch {
		FeatureTrackingBranchProgram(trackingBranch, syncStrategy, FeatureTrackingArgs{
			DevRemote:          args.devRemote,
			FirstCommitMessage: args.firstCommitMessage,
			LastRunParentSHA:   args.parentLastRunSHA,
			LocalName:          args.localName,
			Offline:            args.offline,
			Program:            args.program,
			PushBranches:       args.pushBranches,
		})
	}
	if args.prune {
		args.program.Value.Add(&opcodes.BranchDeleteIfEmptyAtRuntime{Branch: args.localName})
	}
}

type featureBranchArgs struct {
	devRemote          gitdomain.Remote
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

func syncFeatureParentBranch(syncStrategy configdomain.SyncStrategy, args featureBranchArgs) {
	switch syncStrategy {
	case configdomain.SyncStrategyMerge:
		args.program.Value.Add(
			&opcodes.MergeParentsUntilLocal{
				Branch:             args.localName,
				OriginalParentName: args.originalParentName,
				OriginalParentSHA:  args.originalParentSHA,
			},
		)
	case configdomain.SyncStrategyRebase:
		args.program.Value.Add(
			&opcodes.RebaseParentsUntilLocal{
				Branch:      args.localName,
				PreviousSHA: args.parentLastRunSHA,
			},
		)
	case configdomain.SyncStrategyCompress:
		args.program.Value.Add(
			&opcodes.MergeParentsUntilLocal{
				Branch:             args.localName,
				OriginalParentName: args.originalParentName,
				OriginalParentSHA:  args.originalParentSHA,
			},
		)
	case configdomain.SyncStrategyFFOnly:
		// The ff-only strategy does not sync with the parent branch.
		// It is intended for perennial branches only.
	}
}

// separate pull and push of the tracking branch here?
func FeatureTrackingBranchProgram(trackingBranch gitdomain.RemoteBranchName, syncStrategy configdomain.SyncStrategy, args FeatureTrackingArgs) {
	switch syncStrategy {
	case configdomain.SyncStrategyCompress:
		args.Program.Value.Add(
			&opcodes.CompressMergeTrackingBranch{
				CurrentBranch:  args.LocalName,
				CommitMessage:  args.FirstCommitMessage,
				Offline:        args.Offline,
				TrackingBranch: trackingBranch,
			},
		)
	case configdomain.SyncStrategyMerge:
		args.Program.Value.Add(&opcodes.MergeIntoCurrentBranch{BranchToMerge: trackingBranch.BranchName()})
	case configdomain.SyncStrategyRebase:
		if args.Offline.IsFalse() {
			args.Program.Value.Add(
				&opcodes.RebaseTrackingBranch{
					RemoteBranch: trackingBranch,
					PushBranches: args.PushBranches,
				},
				&opcodes.RebaseParentsUntilLocal{
					Branch:      args.LocalName,
					PreviousSHA: args.LastRunParentSHA,
				},
				&opcodes.PushCurrentBranchForceIfNeeded{
					CurrentBranch:   args.LocalName,
					ForceIfIncludes: true,
				},
			)
		}
	case configdomain.SyncStrategyFFOnly:
		if args.Offline.IsFalse() {
			args.Program.Value.Add(&opcodes.MergeFastForward{Branch: trackingBranch.BranchName()})
		}
	}
}

type FeatureTrackingArgs struct {
	DevRemote          gitdomain.Remote
	FirstCommitMessage Option[gitdomain.CommitMessage]
	LastRunParentSHA   Option[gitdomain.SHA]
	LocalName          gitdomain.LocalBranchName
	Offline            configdomain.Offline     // whether offline mode is enabled
	Program            Mutable[program.Program] // the program to update
	PushBranches       configdomain.PushBranches
}
