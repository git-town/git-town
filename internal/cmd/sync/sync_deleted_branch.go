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
	list.Value.Add(&opcodes.CheckoutIfNeeded{Branch: branch})
	pullParentBranchOfCurrentFeatureBranchOpcode(pullParentBranchOfCurrentFeatureBranchOpcodeArgs{
		branch:       branch,
		program:      list,
		syncStrategy: args.Config.SyncFeatureStrategy,
	})
	list.Value.Add(&opcodes.BranchDeleteIfEmptyAtRuntime{Branch: branch})
}

// deletes the given local branch as part of syncing it
func syncDeleteLocalBranchProgram(list Mutable[program.Program], branch gitdomain.LocalBranchName, args BranchProgramArgs) {
	parent := args.Config.Lineage.Parent(branch).GetOrElse(args.Config.MainBranch)
	RemoveBranchConfiguration(RemoveBranchConfigurationArgs{
		Branch:  branch,
		Lineage: args.Config.Lineage,
		Program: list,
	})
	list.Value.Add(&opcodes.CheckoutIfNeeded{Branch: parent})
	list.Value.Add(&opcodes.BranchLocalDelete{Branch: branch})
	list.Value.Add(&opcodes.MessageQueue{Message: fmt.Sprintf(messages.BranchDeleted, branch)})
}
