package sync

import (
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/vm/opcodes"
	"github.com/git-town/git-town/v15/internal/vm/program"
	. "github.com/git-town/git-town/v15/pkg/prelude"
)

// FeatureBranchProgram adds the opcodes to sync the feature branch with the given name.
func FeatureBranchProgram(args featureBranchArgs) {
	syncArgs := syncFeatureBranchProgramArgs{
		firstCommitMessage:  args.firstCommitMessage,
		localName:           args.localName,
		offline:             args.offline,
		parent:              args.parent,
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
	case configdomain.SyncStrategyCompress:
		syncFeatureBranchCompressProgram(syncArgs)
	}
}

type featureBranchArgs struct {
	firstCommitMessage  Option[gitdomain.CommitMessage]
	localName           gitdomain.LocalBranchName
	offline             configdomain.Offline // whether offline mode is enabled
	parent              gitdomain.BranchName
	parentOtherWorktree bool                     // whether the parent of this branch exists on another worktre
	program             Mutable[program.Program] // the program to update
	pushBranches        configdomain.PushBranches
	remoteName          Option[gitdomain.RemoteBranchName]
	syncStrategy        configdomain.SyncStrategy // the sync-feature-strategy
}

func syncFeatureBranchCompressProgram(args syncFeatureBranchProgramArgs) {
	trackingBranch, hasTrackingBranch := args.remoteName.Get()
	if hasTrackingBranch {
		args.program.Value.Add(&opcodes.Merge{Branch: trackingBranch.BranchName()})
	}
	args.program.Value.Add(&opcodes.MergeParent{CurrentBranch: args.localName, ParentActiveInOtherWorktree: args.parentOtherWorktree})
	if firstCommitMessage, branchHasCommits := args.firstCommitMessage.Get(); branchHasCommits {
		args.program.Value.Add(&opcodes.ResetCurrentBranchToParent{CurrentBranch: args.localName})
		args.program.Value.Add(&opcodes.CommitOpenChanges{Message: firstCommitMessage})
	}
	if hasTrackingBranch && args.offline.IsFalse() {
		args.program.Value.Add(&opcodes.ForcePushCurrentBranch{ForceIfIncludes: false})
	}
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
	firstCommitMessage  Option[gitdomain.CommitMessage]
	localName           gitdomain.LocalBranchName
	offline             configdomain.Offline // whether offline mode is enabled
	parent              gitdomain.BranchName
	parentOtherWorktree bool
	program             Mutable[program.Program]
	pushBranches        configdomain.PushBranches
	remoteName          Option[gitdomain.RemoteBranchName]
}
