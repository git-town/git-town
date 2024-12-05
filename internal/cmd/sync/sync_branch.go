package sync

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/config"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/opcodes"
	"github.com/git-town/git-town/v16/internal/vm/program"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/git-town/git-town/v16/pkg/set"
)

// BranchProgram syncs the given branch.
func BranchProgram(localName gitdomain.LocalBranchName, branchInfo gitdomain.BranchInfo, firstCommitMessage Option[gitdomain.CommitMessage], args Mutable[BranchProgramArgs]) {
	originalParentName := args.Value.Config.NormalConfig.Lineage.Parent(localName)
	originalParentSHA := None[gitdomain.SHA]()
	parentName, hasParentName := originalParentName.Get()
	if hasParentName {
		if parentBranchInfo, hasParentBranchInfo := args.Value.BranchInfos.FindLocalOrRemote(parentName).Get(); hasParentBranchInfo {
			originalParentSHA = parentBranchInfo.LocalSHA.Or(parentBranchInfo.RemoteSHA)
		}
	}
	trackingBranchIsGone := branchInfo.SyncStatus == gitdomain.SyncStatusDeletedAtRemote
	rebaseSyncStrategy := args.Value.Config.NormalConfig.SyncFeatureStrategy == configdomain.SyncFeatureStrategyRebase
	hasDescendents := args.Value.Config.NormalConfig.Lineage.HasDescendents(localName)
	parentBranchInfo, hasParentBranchInfo := args.Value.BranchInfos.FindByLocalName(parentName).Get()
	parentTrackingBranchIsGone := false
	if hasParentBranchInfo {
		parentTrackingBranchIsGone = parentBranchInfo.SyncStatus == gitdomain.SyncStatusDeletedAtRemote
	}
	shouldDeleteParent := hasParentName && args.Value.BranchesToDelete.Contains(parentName)
	fmt.Println("1111111111111111111111111111111 branch to sync", localName)
	fmt.Println("1111111111111111111111111111111 BranchesToDelete, shouldDeleteParent", args.Value.BranchesToDelete, shouldDeleteParent)
	fmt.Println("1111111111111111111111111111111 hasParentName, parentName", hasParentName, parentName)
	fmt.Println("1111111111111111111111111111111 trackingBranchIsGone", trackingBranchIsGone)
	fmt.Println("1111111111111111111111111111111 parentTrackingBranchIsGone", parentTrackingBranchIsGone)
	fmt.Println("1111111111111111111111111111111 hasDescendents", hasDescendents)
	// TODO: add an E2E test where a branch has two child branches, and then the branch gets shipped at origin
	switch {
	case rebaseSyncStrategy && trackingBranchIsGone && hasDescendents && shouldDeleteParent:
		// This branch needs to be deleted, and its parent also needs to be deleted.
		args.Value.BranchesToDelete.Add(localName)
	case rebaseSyncStrategy && trackingBranchIsGone && hasDescendents:
		fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		if shouldDeleteParent {
			fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
			removeParentCommits(args.Value.Program, localName, parentName.BranchName(), args.Value.Config.ValidatedConfigData.MainBranch)
		}
		// This branch needs to be deleted and its commits removed from all descendent branches.
		// To do that, we mark it to be deleted here, sync its descendents to remove the commits of this branch,
		// and when that is done we will delete this branch.
		// More info at https://github.com/git-town/git-town/issues/4189.
		args.Value.BranchesToDelete.Add(localName)
	case rebaseSyncStrategy && hasParentToDelete && parentToDeleteName != parentName:
		fmt.Println("CCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC")
		// We sync a branch that has a different parent than the branch currently marked to be deleted.
		// It's time to delete the marked branch now since all its descendents have been synced.
		args.Value.Program.Value.Add(
			&opcodes.BranchLocalDelete{Branch: parentToDeleteName},
			&opcodes.LineageBranchRemove{Branch: parentToDeleteName},
		)
		args.Value.BranchesToDelete = None[gitdomain.LocalBranchName]()
	case trackingBranchIsGone:
		fmt.Println("DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD")
		deletedBranchProgram(args.Value.Program, localName, originalParentName, originalParentSHA, *args.Value)
	case branchInfo.SyncStatus == gitdomain.SyncStatusOtherWorktree:
		// cannot sync branches that are active in another worktree
	default:
		fmt.Println("EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE")
		LocalBranchProgram(localName, branchInfo, originalParentName, originalParentSHA, firstCommitMessage, *args.Value)
	}
	args.Value.Program.Value.Add(&opcodes.ProgramEndOfBranch{})
	fmt.Println("444444444444444444444444444444444444444444444444444444444", args.Value.Program)
}

type BranchProgramArgs struct {
	BranchInfos         gitdomain.BranchInfos // the initial BranchInfos, after "git fetch" ran
	Config              config.ValidatedConfig
	InitialBranch       gitdomain.LocalBranchName
	BranchesToDelete    set.Set[gitdomain.LocalBranchName] // branches that should be deleted after the branches are all synced
	PrefetchBranchInfos gitdomain.BranchInfos              // BranchInfos before "git fetch" ran
	Program             Mutable[program.Program]
	PushBranches        configdomain.PushBranches
	Remotes             gitdomain.Remotes
}

// LocalBranchProgram provides the program to sync a local branch.
func LocalBranchProgram(localName gitdomain.LocalBranchName, branchInfo gitdomain.BranchInfo, originalParentName Option[gitdomain.LocalBranchName], originalParentSHA Option[gitdomain.SHA], firstCommitMessage Option[gitdomain.CommitMessage], args BranchProgramArgs) {
	isMainOrPerennialBranch := args.Config.IsMainOrPerennialBranch(localName)
	if isMainOrPerennialBranch && !args.Remotes.HasOrigin() {
		// perennial branch but no remote --> this branch cannot be synced
		return
	}
	args.Program.Value.Add(&opcodes.CheckoutIfNeeded{Branch: localName})
	branchType := args.Config.BranchType(localName)
	switch branchType {
	case configdomain.BranchTypeFeatureBranch:
		FeatureBranchProgram(args.Config.NormalConfig.SyncFeatureStrategy.SyncStrategy(), featureBranchArgs{
			firstCommitMessage: firstCommitMessage,
			localName:          localName,
			offline:            args.Config.NormalConfig.Offline,
			originalParentName: originalParentName,
			originalParentSHA:  originalParentSHA,
			program:            args.Program,
			pushBranches:       args.PushBranches,
			trackingBranchName: branchInfo.RemoteName,
		})
	case
		configdomain.BranchTypePerennialBranch,
		configdomain.BranchTypeMainBranch:
		PerennialBranchProgram(branchInfo, args)
	case configdomain.BranchTypeParkedBranch:
		ParkedBranchProgram(args.Config.NormalConfig.SyncFeatureStrategy.SyncStrategy(), args.InitialBranch, featureBranchArgs{
			firstCommitMessage: firstCommitMessage,
			localName:          localName,
			offline:            args.Config.NormalConfig.Offline,
			originalParentName: originalParentName,
			originalParentSHA:  originalParentSHA,
			program:            args.Program,
			pushBranches:       args.PushBranches,
			trackingBranchName: branchInfo.RemoteName,
		})
	case configdomain.BranchTypeContributionBranch:
		ContributionBranchProgram(args.Program, branchInfo)
	case configdomain.BranchTypeObservedBranch:
		ObservedBranchProgram(branchInfo.RemoteName, args.Program)
	case configdomain.BranchTypePrototypeBranch:
		FeatureBranchProgram(args.Config.NormalConfig.SyncPrototypeStrategy.SyncStrategy(), featureBranchArgs{
			firstCommitMessage: firstCommitMessage,
			localName:          localName,
			offline:            args.Config.NormalConfig.Offline,
			originalParentName: originalParentName,
			originalParentSHA:  originalParentSHA,
			program:            args.Program,
			pushBranches:       false,
			trackingBranchName: branchInfo.RemoteName,
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

func removeParentCommits(program Mutable[program.Program], branch gitdomain.LocalBranchName, parent gitdomain.BranchName, rebaseOnto gitdomain.LocalBranchName) {
	program.Value.Add(
		&opcodes.CheckoutIfNeeded{Branch: branch},
		&opcodes.PullCurrentBranch{},
		&opcodes.RebaseOnto{
			BranchToRebaseAgainst: parent,
			BranchToRebaseOnto:    rebaseOnto,
		},
		&opcodes.PushCurrentBranchForceIfNeeded{ForceIfIncludes: false},
	)
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
