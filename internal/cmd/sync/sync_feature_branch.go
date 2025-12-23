package sync

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// FeatureBranchProgram adds the opcodes to sync the feature branch with the given name.
func FeatureBranchProgram(syncStrategy configdomain.SyncStrategy, args featureBranchArgs) {
	switch syncStrategy {
	case configdomain.SyncStrategyCompress:
		args.program.Value.Add(
			&opcodes.SyncFeatureBranchCompress{
				CurrentBranch:     args.localName,
				CommitMessage:     args.firstCommitMessage,
				Offline:           args.offline,
				InitialParentName: args.initialParentName,
				InitialParentSHA:  args.initialParentSHA,
				PushBranches:      args.pushBranches,
				TrackingBranch:    args.trackingBranch,
			},
		)
	case configdomain.SyncStrategyFFOnly:
		// The ff-only strategy does not sync with the parent branch.
		// It is intended for perennial branches only.
		if args.offline.IsOnline() {
			if trackingBranch, hasTrackingBranch := args.trackingBranch.Get(); hasTrackingBranch {
				args.program.Value.Add(&opcodes.MergeFastForward{Branch: trackingBranch.BranchName()})
			}
		}
	case configdomain.SyncStrategyMerge:
		args.program.Value.Add(
			&opcodes.SyncFeatureBranchMerge{
				Branch:            args.localName,
				InitialParentName: args.initialParentName,
				InitialParentSHA:  args.initialParentSHA,
				TrackingBranch:    args.trackingBranch,
			},
		)
	case configdomain.SyncStrategyRebase:
		args.program.Value.Add(
			&opcodes.SyncFeatureBranchRebase{
				Branch:               args.localName,
				ParentSHAPreviousRun: args.parentSHAPreviousRun,
				PushBranches:         args.pushBranches,
				TrackingBranch:       args.trackingBranch,
			},
		)
	}
	if args.prune {
		args.program.Value.Add(&opcodes.BranchDeleteIfEmptyAtRuntime{
			Branch:         args.localName,
			TrackingBranch: args.trackingBranch,
		})
	}
}

type featureBranchArgs struct {
	firstCommitMessage   Option[gitdomain.CommitMessage]
	initialParentName    Option[gitdomain.LocalBranchName] // the parent when Git Town started
	initialParentSHA     Option[gitdomain.SHA]             // the parent when Git Town started
	localName            gitdomain.LocalBranchName         // name of the feature branch
	offline              configdomain.Offline              // whether offline mode is enabled
	parentSHAPreviousRun Option[gitdomain.SHA]             // the parent at the end of the last Git Town command
	program              Mutable[program.Program]          // the program to update
	prune                configdomain.Prune
	pushBranches         configdomain.PushBranches
	trackingBranch       Option[gitdomain.RemoteBranchName]
}
