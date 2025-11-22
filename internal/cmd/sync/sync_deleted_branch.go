package sync

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// deletedBranchProgram adds opcodes that sync a branch that was deleted at origin to the given program.
func deletedBranchProgram(branch gitdomain.LocalBranchName, initialParentName Option[gitdomain.LocalBranchName], initialParentSHA, parentSHAPreviousRun Option[gitdomain.SHA], args BranchProgramArgs) {
	switch args.Config.BranchType(branch) {
	case configdomain.BranchTypeFeatureBranch:
		syncDeletedFeatureBranchProgram(branch, initialParentName, initialParentSHA, parentSHAPreviousRun, args)
	case
		configdomain.BranchTypePerennialBranch,
		configdomain.BranchTypeMainBranch,
		configdomain.BranchTypeObservedBranch,
		configdomain.BranchTypeContributionBranch,
		configdomain.BranchTypeParkedBranch,
		configdomain.BranchTypePrototypeBranch:
		syncDeleteLocalBranchProgram(branch, args)
	}
	if _, hasOverride := args.Config.NormalConfig.BranchTypeOverrides[branch]; hasOverride {
		args.Program.Value.Add(&opcodes.BranchTypeOverrideRemove{
			Branch: branch,
		})
	}
}

// syncDeletedFeatureBranchProgram syncs a feare branch whose remote has been deleted.
// The parent branch must have been fully synced before calling this function.
func syncDeletedFeatureBranchProgram(branch gitdomain.LocalBranchName, initialParentName Option[gitdomain.LocalBranchName], initialParentSHA, parentSHAPreviousRun Option[gitdomain.SHA], args BranchProgramArgs) {
	var syncStatus gitdomain.SyncStatus
	if preFetchBranchInfo, has := args.PrefetchBranchInfos.FindByLocalName(branch).Get(); has {
		syncStatus = preFetchBranchInfo.SyncStatus
	} else {
		syncStatus = gitdomain.SyncStatusNotInSync
	}
	switch syncStatus {
	case
		gitdomain.SyncStatusUpToDate,
		gitdomain.SyncStatusBehind,
		gitdomain.SyncStatusLocalOnly:
		syncDeleteLocalBranchProgram(branch, args)
	case
		gitdomain.SyncStatusOtherWorktree,
		gitdomain.SyncStatusRemoteOnly:
		return
	case
		gitdomain.SyncStatusAhead,
		gitdomain.SyncStatusDeletedAtRemote,
		gitdomain.SyncStatusNotInSync:
		args.Program.Value.Add(&opcodes.CheckoutIfNeeded{Branch: branch})
		pullParentBranchOfCurrentFeatureBranchOpcode(pullParentBranchOfCurrentFeatureBranchOpcodeArgs{
			branch:            branch,
			parentNameInitial: initialParentName,
			parentSHAInitial:  initialParentSHA,
			parentSHAPrevious: parentSHAPreviousRun,
			program:           args.Program,
			syncStrategy:      args.Config.NormalConfig.SyncFeatureStrategy,
			// this function syncs a branch whose remote was deleted --> we know for sure there is no tracking branch
			trackingBranch: None[gitdomain.RemoteBranchName](),
		})
		args.Program.Value.Add(&opcodes.BranchWithRemoteGoneDeleteIfEmptyAtRuntime{Branch: branch})
	}
}

// deletes the given local branch as part of syncing it
func syncDeleteLocalBranchProgram(branch gitdomain.LocalBranchName, args BranchProgramArgs) {
	args.Program.Value.Add(
		&opcodes.CheckoutAncestorOrOtherIfNeeded{
			Branch: branch,
		},
		&opcodes.BranchLocalDeleteContent{
			BranchToDelete:     branch,
			BranchToRebaseOnto: args.Config.ValidatedConfigData.MainBranch,
		},
	)
	RemoveBranchConfiguration(RemoveBranchConfigurationArgs{
		Branch:  branch,
		Lineage: args.Config.NormalConfig.Lineage,
		Order:   args.Config.NormalConfig.Order,
		Program: args.Program,
	})
}
