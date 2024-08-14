package sync

import (
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/git-town/git-town/v15/internal/vm/opcodes"
	"github.com/git-town/git-town/v15/internal/vm/program"
)

// FeatureBranchProgram adds the opcodes to sync the feature branch with the given name.
func FeatureBranchProgram(args featureBranchArgs) {
	syncArgs := syncFeatureBranchProgramArgs{
		localName:           args.localName,
		offline:             args.offline,
		parentOtherWorktree: args.parentOtherWorktree,
		program:             args.program,
		pushBranches:        args.pushBranches,
		remoteName:          args.remoteName,
	}
	switch args.syncStrategy {
	case configdomain.SyncStrategyMerge:
		syncFeatureBranchMergeProgram(syncArgs)
	case configdomain.SyncStrategyRebase:
		syncFeatureBranchRebaseProgram(syncArgs)
	}
}

type featureBranchArgs struct {
	localName           gitdomain.LocalBranchName
	offline             configdomain.Offline     // whether offline mode is enabled
	parentOtherWorktree bool                     // whether the parent of this branch exists on another worktre
	program             Mutable[program.Program] // the program to update
	pushBranches        configdomain.PushBranches
	remoteName          Option[gitdomain.RemoteBranchName]
	syncStrategy        configdomain.SyncStrategy // the sync-feature-strategy
}

// syncs the given feature branch using the "merge" sync strategy
func syncFeatureBranchMergeProgram(args syncFeatureBranchProgramArgs) {
	if trackingBranch, hasTrackingBranch := args.remoteName.Get(); hasTrackingBranch {
		args.program.Value.Add(&opcodes.Merge{Branch: trackingBranch.BranchName()})
	}
	args.program.Value.Add(&opcodes.MergeParent{CurrentBranch: args.localName, ParentActiveInOtherWorktree: args.parentOtherWorktree})
}

// syncs the given feature branch using the "rebase" sync strategy
func syncFeatureBranchRebaseProgram(args syncFeatureBranchProgramArgs) {
	// rebase against parent
	args.program.Value.Add(&opcodes.RebaseParent{
		CurrentBranch:               args.localName,
		ParentActiveInOtherWorktree: args.parentOtherWorktree,
	})
	if trackingBranch, hasTrackingBranch := args.remoteName.Get(); hasTrackingBranch {
		if args.offline.IsFalse() {
			args.program.Value.Add(&opcodes.RebaseFeatureTrackingBranch{RemoteBranch: trackingBranch, PushBranches: args.pushBranches})
		}
	}
}

type syncFeatureBranchProgramArgs struct {
	localName           gitdomain.LocalBranchName
	offline             configdomain.Offline // whether offline mode is enabled
	parentOtherWorktree bool
	program             Mutable[program.Program]
	pushBranches        configdomain.PushBranches
	remoteName          Option[gitdomain.RemoteBranchName]
}
