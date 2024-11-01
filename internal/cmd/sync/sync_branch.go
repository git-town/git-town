package sync

import (
	"github.com/git-town/git-town/v16/internal/config"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/opcodes"
	"github.com/git-town/git-town/v16/internal/vm/program"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// BranchProgram syncs the given branch.
func BranchProgram(localName gitdomain.LocalBranchName, branchInfo gitdomain.BranchInfo, firstCommitMessage Option[gitdomain.CommitMessage], args BranchProgramArgs) {
	// TODO: calculate originalParentName and originalParentSHA here and provide as arguments.
	// Also, possibly provide the BranchInfo of the original parent branch here, because both values can be derived from it.
	// Don't store the branchinfo in the opcodes, though.
	switch {
	case branchInfo.SyncStatus == gitdomain.SyncStatusDeletedAtRemote:
		deletedBranchProgram(args.Program, localName, args)
	case branchInfo.SyncStatus == gitdomain.SyncStatusOtherWorktree:
		// Git Town doesn't sync branches that are active in another worktree
	default:
		localBranchProgram(localName, branchInfo, firstCommitMessage, args)
	}
	args.Program.Value.Add(&opcodes.ProgramEndOfBranch{})
}

type BranchProgramArgs struct {
	BranchInfos         gitdomain.BranchInfos // the initial BranchInfos, after "git fetch" ran
	Config              config.ValidatedConfig
	InitialBranch       gitdomain.LocalBranchName
	PrefetchBranchInfos gitdomain.BranchInfos // BranchInfos before "git fetch" ran
	Program             Mutable[program.Program]
	PushBranches        configdomain.PushBranches
	Remotes             gitdomain.Remotes
}

// localBranchProgram provides the program to sync a local branch.
func localBranchProgram(localName gitdomain.LocalBranchName, branchInfo gitdomain.BranchInfo, firstCommitMessage Option[gitdomain.CommitMessage], args BranchProgramArgs) {
	isMainOrPerennialBranch := args.Config.IsMainOrPerennialBranch(localName)
	if isMainOrPerennialBranch && !args.Remotes.HasOrigin() {
		// perennial branch but no remote --> this branch cannot be synced
		return
	}
	args.Program.Value.Add(&opcodes.CheckoutIfNeeded{Branch: localName})
	branchType := args.Config.BranchType(localName)
	parentNameOpt := args.Config.NormalConfig.Lineage.Parent(localName)
	parentSHAOpt := None[gitdomain.SHA]()
	if parentName, hasParentName := parentNameOpt.Get(); hasParentName {
		if parentBranchInfo, hasParentBranchInfo := args.BranchInfos.FindLocalOrRemote(parentName).Get(); hasParentBranchInfo {
			parentSHAOpt = parentBranchInfo.LocalSHA.Or(parentBranchInfo.RemoteSHA)
		}
	}
	switch branchType {
	case configdomain.BranchTypeFeatureBranch:
		FeatureBranchProgram(featureBranchArgs{
			firstCommitMessage: firstCommitMessage,
			localName:          localName,
			offline:            args.Config.NormalConfig.Offline,
			originalParentName: parentNameOpt,
			originalParentSHA:  parentSHAOpt,
			program:            args.Program,
			pushBranches:       args.PushBranches,
			remoteName:         branchInfo.RemoteName,
			syncStrategy:       args.Config.NormalConfig.SyncFeatureStrategy.SyncStrategy(),
		})
	case
		configdomain.BranchTypePerennialBranch,
		configdomain.BranchTypeMainBranch:
		PerennialBranchProgram(branchInfo, args)
	case configdomain.BranchTypeParkedBranch:
		ParkedBranchProgram(args.InitialBranch, featureBranchArgs{
			firstCommitMessage: firstCommitMessage,
			localName:          localName,
			offline:            args.Config.NormalConfig.Offline,
			originalParentName: parentNameOpt,
			originalParentSHA:  parentSHAOpt,
			program:            args.Program,
			pushBranches:       args.PushBranches,
			remoteName:         branchInfo.RemoteName,
			syncStrategy:       args.Config.NormalConfig.SyncFeatureStrategy.SyncStrategy(),
		})
	case configdomain.BranchTypeContributionBranch:
		ContributionBranchProgram(args.Program, branchInfo)
	case configdomain.BranchTypeObservedBranch:
		ObservedBranchProgram(branchInfo.RemoteName, args.Program)
	case configdomain.BranchTypePrototypeBranch:
		FeatureBranchProgram(featureBranchArgs{
			firstCommitMessage: firstCommitMessage,
			localName:          localName,
			offline:            args.Config.NormalConfig.Offline,
			originalParentName: parentNameOpt,
			originalParentSHA:  parentSHAOpt,
			program:            args.Program,
			pushBranches:       false,
			remoteName:         branchInfo.RemoteName,
			syncStrategy:       args.Config.NormalConfig.SyncPrototypeStrategy.SyncStrategy(),
		})
	}
	if args.PushBranches.IsTrue() && args.Remotes.HasOrigin() && args.Config.NormalConfig.IsOnline() && branchType.ShouldPush(localName == args.InitialBranch) {
		switch {
		case !branchInfo.HasTrackingBranch():
			args.Program.Value.Add(&opcodes.BranchTrackingCreate{Branch: localName})
		case isMainOrPerennialBranch:
			args.Program.Value.Add(&opcodes.PushCurrentBranchIfNeeded{CurrentBranch: localName})
		default:
			pushFeatureBranchProgram(args.Program, localName, args.Config.NormalConfig.SyncFeatureStrategy)
		}
	}
}

// pullParentBranchOfCurrentFeatureBranchOpcode adds the opcode to pull updates from the parent branch of the current feature branch into the current feature branch.
func pullParentBranchOfCurrentFeatureBranchOpcode(args pullParentBranchOfCurrentFeatureBranchOpcodeArgs) {
	switch args.syncStrategy {
	case configdomain.SyncFeatureStrategyMerge:
		args.program.Value.Add(&opcodes.MergeParentIfNeeded{
			Branch:             args.branch,
			OriginalParentName: args.originalParentName,
			OriginalParentSHA:  args.originalParentSHA,
		})
	case configdomain.SyncFeatureStrategyRebase:
		args.program.Value.Add(&opcodes.RebaseParentIfNeeded{
			Branch: args.branch,
		})
	case configdomain.SyncFeatureStrategyCompress:
		args.program.Value.Add(&opcodes.MergeParentIfNeeded{
			Branch:             args.branch,
			OriginalParentName: args.originalParentName,
			OriginalParentSHA:  args.originalParentSHA,
		})
	}
}

type pullParentBranchOfCurrentFeatureBranchOpcodeArgs struct {
	branch             gitdomain.LocalBranchName
	originalParentName Option[gitdomain.LocalBranchName]
	originalParentSHA  Option[gitdomain.SHA]
	program            Mutable[program.Program]
	syncStrategy       configdomain.SyncFeatureStrategy
}

func pushFeatureBranchProgram(prog Mutable[program.Program], branch gitdomain.LocalBranchName, syncFeatureStrategy configdomain.SyncFeatureStrategy) {
	switch syncFeatureStrategy {
	case configdomain.SyncFeatureStrategyMerge:
		prog.Value.Add(&opcodes.PushCurrentBranchIfNeeded{CurrentBranch: branch})
	case configdomain.SyncFeatureStrategyRebase:
		prog.Value.Add(&opcodes.PushCurrentBranchForceIfNeeded{ForceIfIncludes: true})
	case configdomain.SyncFeatureStrategyCompress:
		prog.Value.Add(&opcodes.PushCurrentBranchForceIfNeeded{ForceIfIncludes: false})
	}
}

// updateCurrentPerennialBranchOpcode provides the opcode to update the current perennial branch with changes from the given other branch.
func updateCurrentPerennialBranchOpcode(prog Mutable[program.Program], otherBranch gitdomain.RemoteBranchName, strategy configdomain.SyncPerennialStrategy) {
	switch strategy {
	case configdomain.SyncPerennialStrategyMerge:
		prog.Value.Add(&opcodes.Merge{Branch: otherBranch.BranchName()})
	case configdomain.SyncPerennialStrategyRebase:
		prog.Value.Add(&opcodes.RebaseBranch{Branch: otherBranch.BranchName()})
	}
}
