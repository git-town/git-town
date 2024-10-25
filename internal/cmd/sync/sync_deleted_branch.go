package sync

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/vm/opcodes"
	"github.com/git-town/git-town/v16/internal/vm/program"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// deletedBranchProgram adds opcodes that sync a branch that was deleted at origin to the given program.
func deletedBranchProgram(prog Mutable[program.Program], branch gitdomain.LocalBranchName, args BranchProgramArgs) {
	switch args.Config.BranchType(branch) {
	case configdomain.BranchTypeFeatureBranch:
		syncDeletedFeatureBranchProgram(prog, branch, args)
	case
		configdomain.BranchTypePerennialBranch,
		configdomain.BranchTypeMainBranch,
		configdomain.BranchTypeObservedBranch,
		configdomain.BranchTypeContributionBranch,
		configdomain.BranchTypeParkedBranch,
		configdomain.BranchTypePrototypeBranch:
		syncDeleteLocalBranchProgram(prog, branch, args)
	}
}

// syncDeletedFeatureBranchProgram syncs a feare branch whose remote has been deleted.
// The parent branch must have been fully synced before calling this function.
func syncDeletedFeatureBranchProgram(prog Mutable[program.Program], branch gitdomain.LocalBranchName, args BranchProgramArgs) {
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
		syncDeleteLocalBranchProgram(prog, branch, args)
	case gitdomain.SyncStatusOtherWorktree:
	case gitdomain.SyncStatusRemoteOnly:
		return
	case gitdomain.SyncStatusAhead:
	case gitdomain.SyncStatusDeletedAtRemote:
	case gitdomain.SyncStatusNotInSync:
		prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: branch})
		pullParentBranchOfCurrentFeatureBranchOpcode(pullParentBranchOfCurrentFeatureBranchOpcodeArgs{
			branch:       branch,
			program:      prog,
			syncStrategy: args.Config.NormalConfig.SyncFeatureStrategy,
		})
		prog.Value.Add(&opcodes.BranchDeleteIfEmptyAtRuntime{Branch: branch})
	}
}

// deletes the given local branch as part of syncing it
func syncDeleteLocalBranchProgram(prog Mutable[program.Program], branch gitdomain.LocalBranchName, args BranchProgramArgs) {
	RemoveBranchConfiguration(RemoveBranchConfigurationArgs{
		Branch:  branch,
		Lineage: args.Config.NormalConfig.Lineage,
		Program: prog,
	})
	prog.Value.Add(&opcodes.CheckoutParentOrMain{Branch: branch})
	prog.Value.Add(&opcodes.BranchLocalDelete{Branch: branch})
	prog.Value.Add(&opcodes.MessageQueue{Message: fmt.Sprintf(messages.BranchDeleted, branch)})
}
