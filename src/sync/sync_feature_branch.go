package sync

import (
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/vm/opcodes"
	"github.com/git-town/git-town/v14/src/vm/program"
)

// FeatureBranchProgram adds the opcodes to sync the feature branch with the given name.
func FeatureBranchProgram(args featureBranchArgs) {
	syncArgs := syncFeatureBranchProgramArgs{
		localName:           args.localName,
		offline:             args.offline,
		parentOtherWorktree: args.parentOtherWorktree,
		program:             args.program,
		remoteName:          args.remoteName,
	}
	switch args.syncStrategy {
	case configdomain.SyncFeatureStrategyMerge:
		syncFeatureBranchMergeProgram(syncArgs)
	case configdomain.SyncFeatureStrategyRebase:
		syncFeatureBranchRebaseProgram(syncArgs)
	}
}

type featureBranchArgs struct {
	localName           gitdomain.LocalBranchName
	offline             configdomain.Offline // whether offline mode is enabled
	parentOtherWorktree bool                 // whether the parent of this branch exists on another worktre
	program             *program.Program     // the program to update
	remoteName          Option[gitdomain.RemoteBranchName]
	syncStrategy        configdomain.SyncFeatureStrategy // the sync-feature-strategy
}

// syncs the given feature branch using the "merge" sync strategy
func syncFeatureBranchMergeProgram(args syncFeatureBranchProgramArgs) {
	if trackingBranch, hasTrackingBranch := args.remoteName.Get(); hasTrackingBranch {
		args.program.Add(&opcodes.Merge{Branch: trackingBranch.BranchName()})
	}
	args.program.Add(&opcodes.MergeParent{CurrentBranch: args.localName, ParentActiveInOtherWorktree: args.parentOtherWorktree})
}

// syncs the given feature branch using the "rebase" sync strategy
func syncFeatureBranchRebaseProgram(args syncFeatureBranchProgramArgs) {
	// rebase against parent
	args.program.Add(&opcodes.RebaseParent{
		CurrentBranch:               args.localName,
		ParentActiveInOtherWorktree: args.parentOtherWorktree,
	})
	if trackingBranch, hasTrackingBranch := args.remoteName.Get(); hasTrackingBranch {
		if !args.offline.Bool() {
			args.program.Add(&opcodes.RebaseFeatureTrackingBranch{RemoteBranch: trackingBranch})
		}
	}
}

type syncFeatureBranchProgramArgs struct {
	localName           gitdomain.LocalBranchName
	offline             configdomain.Offline // whether offline mode is enabled
	parentOtherWorktree bool
	program             *program.Program
	remoteName          Option[gitdomain.RemoteBranchName]
}
