package sync

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/vm/opcodes"
	"github.com/git-town/git-town/v14/src/vm/program"
)

// syncDeletedBranchProgram adds opcodes that sync a branch that was deleted at origin to the given program.
func syncDeletedBranchProgram(list Mutable[program.Program], branch gitdomain.LocalBranchName, parentOtherWorktree bool, args BranchProgramArgs) {
	switch args.Config.BranchType(branch) {
	case configdomain.BranchTypeFeatureBranch:
		syncDeletedFeatureBranchProgram(list, branch, parentOtherWorktree, args)
	case configdomain.BranchTypePerennialBranch, configdomain.BranchTypeMainBranch:
		syncDeletedPerennialBranchProgram(list, branch, args)
	case configdomain.BranchTypeObservedBranch, configdomain.BranchTypeContributionBranch, configdomain.BranchTypeParkedBranch:
		syncDeletedObservedBranchProgram(list, branch, args)
	case configdomain.BranchTypePrototypeBranch:
		syncDeletedPrototypeBranchProgram(list, branch, args)
	}
}

// syncDeletedFeatureBranchProgram syncs a feare branch whose remote has been deleted.
// The parent branch must have been fully synced before calling this function.
func syncDeletedFeatureBranchProgram(list Mutable[program.Program], branch gitdomain.LocalBranchName, parentOtherWorktree bool, args BranchProgramArgs) {
	list.Value.Add(&opcodes.Checkout{Branch: branch})
	pullParentBranchOfCurrentFeatureBranchOpcode(pullParentBranchOfCurrentFeatureBranchOpcodeArgs{
		branch:              branch,
		parentOtherWorktree: parentOtherWorktree,
		program:             list,
		syncStrategy:        args.Config.SyncFeatureStrategy,
	})
	list.Value.Add(&opcodes.DeleteBranchIfEmptyAtRuntime{Branch: branch})
}

func syncDeletedObservedBranchProgram(list Mutable[program.Program], branch gitdomain.LocalBranchName, args BranchProgramArgs) {
	RemoveBranchFromLineage(RemoveBranchFromLineageArgs{
		Branch:  branch,
		Lineage: args.Config.Lineage,
		Parent:  args.Config.MainBranch, // TODO: change from MainBranch to the actual parent of the branch whose remote was deleted
		Program: list,
	})
	list.Value.Add(&opcodes.RemoveFromObservedBranches{Branch: branch})
	list.Value.Add(&opcodes.Checkout{Branch: args.Config.MainBranch}) // TODO: checkout the actual parent here
	list.Value.Add(&opcodes.DeleteLocalBranch{Branch: branch})
	list.Value.Add(&opcodes.QueueMessage{Message: fmt.Sprintf(messages.BranchDeleted, branch)})
}

func syncDeletedPerennialBranchProgram(list Mutable[program.Program], branch gitdomain.LocalBranchName, args BranchProgramArgs) {
	RemoveBranchFromLineage(RemoveBranchFromLineageArgs{
		Branch:  branch,
		Lineage: args.Config.Lineage,
		Parent:  args.Config.MainBranch, // TODO: change from MainBranch to the actual parent of the branch whose remote was deleted
		Program: list,
	})
	list.Value.Add(&opcodes.RemoveFromPerennialBranches{Branch: branch})
	list.Value.Add(&opcodes.Checkout{Branch: args.Config.MainBranch}) // TODO: checkout the actual parent here
	list.Value.Add(&opcodes.DeleteLocalBranch{Branch: branch})
	list.Value.Add(&opcodes.QueueMessage{Message: fmt.Sprintf(messages.BranchDeleted, branch)})
}

func syncDeletedPrototypeBranchProgram(list Mutable[program.Program], branch gitdomain.LocalBranchName, args BranchProgramArgs) {
	RemoveBranchFromLineage(RemoveBranchFromLineageArgs{
		Branch:  branch,
		Lineage: args.Config.Lineage,
		Parent:  args.Config.MainBranch, // TODO: change from MainBranch to the actual parent of the branch whose remote was deleted
		Program: list,
	})
	list.Value.Add(&opcodes.RemoveFromPrototypeBranches{Branch: branch})
	list.Value.Add(&opcodes.Checkout{Branch: args.Config.MainBranch}) // TODO: checkout the actual parent here
	list.Value.Add(&opcodes.DeleteLocalBranch{Branch: branch})
	list.Value.Add(&opcodes.QueueMessage{Message: fmt.Sprintf(messages.BranchDeleted, branch)})
}
