package sync

import (
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/opcodes"
	"github.com/git-town/git-town/v21/internal/vm/program"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/git-town/git-town/v21/pkg/set"
)

// BranchProgram syncs the given branch.
func BranchProgram(localName gitdomain.LocalBranchName, branchInfo gitdomain.BranchInfo, firstCommitMessage Option[gitdomain.CommitMessage], args BranchProgramArgs) {
	initialParentName := args.Config.NormalConfig.Lineage.Parent(localName)
	initialParentSHA := None[gitdomain.SHA]()
	parentName, hasParentName := initialParentName.Get()
	if hasParentName {
		if parentBranchInfo, hasParentBranchInfo := args.BranchInfos.FindLocalOrRemote(parentName, args.Config.NormalConfig.DevRemote).Get(); hasParentBranchInfo {
			initialParentSHA = parentBranchInfo.LocalSHA.Or(parentBranchInfo.RemoteSHA)
		}
	}
	trackingBranchGone := branchInfo.SyncStatus == gitdomain.SyncStatusDeletedAtRemote
	rebaseSyncStrategy := args.Config.NormalConfig.SyncFeatureStrategy == configdomain.SyncFeatureStrategyRebase
	hasDescendents := args.Config.NormalConfig.Lineage.HasDescendents(localName)
	parentToRemove, hasParentToRemove := args.Config.NormalConfig.Lineage.LatestAncestor(localName, args.BranchesToDelete.Value.Values()).Get()
	if hasParentToRemove && rebaseSyncStrategy {
		RemoveAncestorCommits(RemoveAncestorCommitsArgs{
			Ancestor:          parentToRemove.BranchName(),
			Branch:            localName,
			HasTrackingBranch: branchInfo.HasTrackingBranch(),
			Program:           args.Program,
			RebaseOnto:        args.Config.ValidatedConfigData.MainBranch, // TODO: RebaseOnto the latest existing parent, which isn't always main
		})
	}
	parentSHAPreviousRun := None[gitdomain.SHA]()
	if parent, has := initialParentName.Get(); has {
		if branchInfosLastRun, has := args.BranchInfosLastRun.Get(); has {
			if parentInfoLastRun, has := branchInfosLastRun.FindByLocalName(parent).Get(); has {
				parentSHAPreviousRun = Some(parentInfoLastRun.GetLocalOrRemoteSHA())
			}
		}
	}
	switch {
	case hasParentToRemove && parentToRemove == parentName && trackingBranchGone && hasDescendents:
		args.BranchesToDelete.Value.Add(localName)
	case hasParentToRemove && parentToRemove == parentName:
		// nothing to do here, we already synced with the parent
	case rebaseSyncStrategy && trackingBranchGone && hasDescendents:
		args.BranchesToDelete.Value.Add(localName)
	case trackingBranchGone:
		deletedBranchProgram(args.Program, localName, initialParentName, initialParentSHA, parentSHAPreviousRun, args)
	case branchInfo.SyncStatus == gitdomain.SyncStatusOtherWorktree:
		// cannot sync branches that are active in another worktree
	default:
		LocalBranchProgram(localName, branchInfo, initialParentName, initialParentSHA, parentSHAPreviousRun, firstCommitMessage, args)
	}
	args.Program.Value.Add(&opcodes.ProgramEndOfBranch{})
}

type BranchProgramArgs struct {
	BranchInfos         gitdomain.BranchInfos                       // the initial BranchInfos, after "git fetch" ran
	BranchInfosLastRun  Option[gitdomain.BranchInfos]               // the BranchInfos at the end of the previous Git Town command
	BranchesToDelete    Mutable[set.Set[gitdomain.LocalBranchName]] // branches that should be deleted after the branches are all synced
	Config              config.ValidatedConfig
	InitialBranch       gitdomain.LocalBranchName
	PrefetchBranchInfos gitdomain.BranchInfos // BranchInfos before "git fetch" ran
	Program             Mutable[program.Program]
	Prune               configdomain.Prune
	PushBranches        configdomain.PushBranches
	Remotes             gitdomain.Remotes
}

// LocalBranchProgram provides the program to sync a local branch.
func LocalBranchProgram(localName gitdomain.LocalBranchName, branchInfo gitdomain.BranchInfo, initialParentName Option[gitdomain.LocalBranchName], initialParentSHA, parentSHAPreviousRun Option[gitdomain.SHA], firstCommitMessage Option[gitdomain.CommitMessage], args BranchProgramArgs) {
	branchType := args.Config.BranchType(localName)
	isMainOrPerennialBranch := branchType == configdomain.BranchTypeMainBranch || branchType == configdomain.BranchTypePerennialBranch
	if isMainOrPerennialBranch && !args.Remotes.HasRemote(args.Config.NormalConfig.DevRemote) {
		// perennial branch but no remote --> this branch cannot be synced
		return
	}
	args.Program.Value.Add(&opcodes.CheckoutIfNeeded{Branch: localName})
	switch branchType {
	case configdomain.BranchTypeFeatureBranch:
		FeatureBranchProgram(args.Config.NormalConfig.SyncFeatureStrategy.SyncStrategy(), featureBranchArgs{
			firstCommitMessage:   firstCommitMessage,
			initialParentName:    initialParentName,
			initialParentSHA:     initialParentSHA,
			localName:            localName,
			offline:              args.Config.NormalConfig.Offline,
			parentSHAPreviousRun: parentSHAPreviousRun,
			program:              args.Program,
			prune:                args.Prune,
			pushBranches:         args.PushBranches,
			trackingBranch:       branchInfo.RemoteName,
		})
	case configdomain.BranchTypePerennialBranch, configdomain.BranchTypeMainBranch:
		PerennialBranchProgram(branchInfo, args)
	case configdomain.BranchTypeParkedBranch:
		ParkedBranchProgram(args.Config.NormalConfig.SyncFeatureStrategy.SyncStrategy(), args.InitialBranch, featureBranchArgs{
			firstCommitMessage:   firstCommitMessage,
			initialParentName:    initialParentName,
			initialParentSHA:     initialParentSHA,
			localName:            localName,
			offline:              args.Config.NormalConfig.Offline,
			parentSHAPreviousRun: parentSHAPreviousRun,
			program:              args.Program,
			prune:                args.Prune,
			pushBranches:         args.PushBranches,
			trackingBranch:       branchInfo.RemoteName,
		})
	case configdomain.BranchTypeContributionBranch:
		ContributionBranchProgram(args.Program, branchInfo)
	case configdomain.BranchTypeObservedBranch:
		ObservedBranchProgram(branchInfo, args.Program)
	case configdomain.BranchTypePrototypeBranch:
		FeatureBranchProgram(args.Config.NormalConfig.SyncPrototypeStrategy.SyncStrategy(), featureBranchArgs{
			firstCommitMessage:   firstCommitMessage,
			initialParentName:    initialParentName,
			initialParentSHA:     initialParentSHA,
			localName:            localName,
			offline:              args.Config.NormalConfig.Offline,
			parentSHAPreviousRun: parentSHAPreviousRun,
			program:              args.Program,
			prune:                args.Prune,
			pushBranches:         configdomain.PushBranches(branchInfo.HasTrackingBranch()),
			trackingBranch:       branchInfo.RemoteName,
		})
	}
	if args.PushBranches.IsTrue() && args.Remotes.HasRemote(args.Config.NormalConfig.DevRemote) && args.Config.NormalConfig.Offline.IsOnline() && branchType.ShouldPush(localName == args.InitialBranch) {
		isMainBranch := branchType == configdomain.BranchTypeMainBranch
		switch {
		case !branchInfo.HasTrackingBranch():
			args.Program.Value.Add(&opcodes.BranchTrackingCreate{Branch: localName})
		case isMainBranch && args.Remotes.HasUpstream() && args.Config.NormalConfig.SyncUpstream.IsTrue():
			args.Program.Value.Add(&opcodes.PushCurrentBranchIfNeeded{CurrentBranch: localName})
		case isMainOrPerennialBranch && !shouldPushPerennialBranch(branchInfo.SyncStatus):
			// don't push if its a perennial branch that doesn't need pushing
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
	case configdomain.SyncFeatureStrategyMerge, configdomain.SyncFeatureStrategyCompress:
		args.program.Value.Add(&opcodes.SyncFeatureBranchMerge{
			Branch:            args.branch,
			InitialParentName: args.initialParentName,
			InitialParentSHA:  args.initialParentSHA,
			TrackingBranch:    args.trackingBranch,
		})
	case configdomain.SyncFeatureStrategyRebase:
		args.program.Value.Add(&opcodes.RebaseParentsUntilLocal{
			Branch:               args.branch,
			ParentSHAPreviousRun: args.parentSHAPreviousRun,
		})
	}
}

type pullParentBranchOfCurrentFeatureBranchOpcodeArgs struct {
	branch               gitdomain.LocalBranchName
	initialParentName    Option[gitdomain.LocalBranchName]
	initialParentSHA     Option[gitdomain.SHA]
	parentSHAPreviousRun Option[gitdomain.SHA]
	program              Mutable[program.Program]
	syncStrategy         configdomain.SyncFeatureStrategy
	trackingBranch       Option[gitdomain.RemoteBranchName]
}

func pushFeatureBranchProgram(prog Mutable[program.Program], branch gitdomain.LocalBranchName, syncFeatureStrategy configdomain.SyncFeatureStrategy) {
	switch syncFeatureStrategy {
	case configdomain.SyncFeatureStrategyMerge:
		prog.Value.Add(&opcodes.PushCurrentBranchIfNeeded{CurrentBranch: branch})
	case configdomain.SyncFeatureStrategyRebase:
		prog.Value.Add(&opcodes.PushCurrentBranchForceIfNeeded{CurrentBranch: branch, ForceIfIncludes: true})
	case configdomain.SyncFeatureStrategyCompress:
		prog.Value.Add(&opcodes.PushCurrentBranchForceIfNeeded{CurrentBranch: branch, ForceIfIncludes: false})
	}
}

func RemoveAncestorCommits(args RemoveAncestorCommitsArgs) {
	args.Program.Value.Add(
		&opcodes.CheckoutIfNeeded{Branch: args.Branch},
	)
	if args.HasTrackingBranch {
		args.Program.Value.Add(
			&opcodes.PullCurrentBranch{},
		)
	}
	args.Program.Value.Add(
		&opcodes.RebaseOntoRemoveDeleted{
			BranchToRebaseOnto: args.RebaseOnto,
			CommitsToRemove:    args.Ancestor,
			Upstream:           None[gitdomain.LocalBranchName](),
		},
	)
	if args.HasTrackingBranch {
		args.Program.Value.Add(
			&opcodes.PushCurrentBranchForce{ForceIfIncludes: false},
		)
	}
}

type RemoveAncestorCommitsArgs struct {
	Ancestor          gitdomain.BranchName
	Branch            gitdomain.LocalBranchName
	HasTrackingBranch bool
	Program           Mutable[program.Program]
	RebaseOnto        gitdomain.LocalBranchName
}

func shouldPushPerennialBranch(syncStatus gitdomain.SyncStatus) bool {
	switch syncStatus {
	case
		gitdomain.SyncStatusAhead,
		gitdomain.SyncStatusBehind,
		gitdomain.SyncStatusLocalOnly,
		gitdomain.SyncStatusNotInSync:
		return true
	case
		gitdomain.SyncStatusDeletedAtRemote,
		gitdomain.SyncStatusOtherWorktree,
		gitdomain.SyncStatusRemoteOnly,
		gitdomain.SyncStatusUpToDate:
	}
	return false
}

// updateCurrentPerennialBranchOpcode provides the opcode to update the current perennial branch with changes from the given other branch.
func updateCurrentPerennialBranchOpcode(prog Mutable[program.Program], otherBranch gitdomain.RemoteBranchName, strategy configdomain.SyncPerennialStrategy) {
	switch strategy {
	case configdomain.SyncPerennialStrategyMerge:
		prog.Value.Add(&opcodes.MergeIntoCurrentBranch{BranchToMerge: otherBranch.BranchName()})
	case configdomain.SyncPerennialStrategyRebase:
		prog.Value.Add(&opcodes.RebaseBranch{Branch: otherBranch.BranchName()})
	case configdomain.SyncPerennialStrategyFFOnly:
		prog.Value.Add(&opcodes.MergeFastForward{Branch: otherBranch.BranchName()})
	}
}
