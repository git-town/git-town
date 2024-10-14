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
func deletedBranchProgram(list Mutable[program.Program], branch gitdomain.LocalBranchName, args BranchProgramArgs) {
	switch args.Config.BranchType(branch) {
	case configdomain.BranchTypeFeatureBranch:
		syncDeletedFeatureBranchProgram(list, branch, args)
	case
		configdomain.BranchTypePerennialBranch,
		configdomain.BranchTypeMainBranch,
		configdomain.BranchTypeObservedBranch,
		configdomain.BranchTypeContributionBranch,
		configdomain.BranchTypeParkedBranch,
		configdomain.BranchTypePrototypeBranch:
		syncDeleteLocalBranchProgram(list, branch, args)
	}
}

// syncDeletedFeatureBranchProgram syncs a feare branch whose remote has been deleted.
// The parent branch must have been fully synced before calling this function.
func syncDeletedFeatureBranchProgram(list Mutable[program.Program], branch gitdomain.LocalBranchName, args BranchProgramArgs) {
	if preFetchBranchInfo, has := args.PrefetchBranchInfos.FindByLocalName(branch).Get(); has {
		switch preFetchBranchInfo.SyncStatus {
		case gitdomain.SyncStatusUpToDate, gitdomain.SyncStatusBehind, gitdomain.SyncStatusLocalOnly:
			list.Value.Add(&opcodes.DeleteLocalBranch{Branch: branch})
		case gitdomain.SyncStatusAhead:
		case gitdomain.SyncStatusDeletedAtRemote:
		case gitdomain.SyncStatusNotInSync:
		case gitdomain.SyncStatusOtherWorktree:
		case gitdomain.SyncStatusRemoteOnly:
		}
	}
	list.Value.Add(&opcodes.CheckoutIfNeeded{Branch: branch})
	pullParentBranchOfCurrentFeatureBranchOpcode(pullParentBranchOfCurrentFeatureBranchOpcodeArgs{
		branch:       branch,
		program:      list,
		syncStrategy: args.Config.SyncFeatureStrategy,
	})
	list.Value.Add(&opcodes.DeleteBranchIfEmptyAtRuntime{Branch: branch})
}

// deletes the given local branch as part of syncing it
func syncDeleteLocalBranchProgram(list Mutable[program.Program], branch gitdomain.LocalBranchName, args BranchProgramArgs) {
	parent := args.Config.Lineage.Parent(branch).GetOrElse(args.Config.MainBranch)
	RemoveBranchFromLineage(RemoveBranchFromLineageArgs{
		Branch:  branch,
		Lineage: args.Config.Lineage,
		Parent:  parent,
		Program: list,
	})
	list.Value.Add(&opcodes.RemoveFromObservedBranches{Branch: branch})
	list.Value.Add(&opcodes.RemoveFromPerennialBranches{Branch: branch})
	list.Value.Add(&opcodes.RemoveFromPrototypeBranches{Branch: branch})
	list.Value.Add(&opcodes.CheckoutIfNeeded{Branch: parent})
	list.Value.Add(&opcodes.DeleteLocalBranch{Branch: branch})
	list.Value.Add(&opcodes.QueueMessage{Message: fmt.Sprintf(messages.BranchDeleted, branch)})
}
