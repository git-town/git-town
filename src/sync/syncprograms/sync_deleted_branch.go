package syncprograms

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/messages"
	"github.com/git-town/git-town/v11/src/undo/undodomain"
	"github.com/git-town/git-town/v11/src/vm/opcode"
	"github.com/git-town/git-town/v11/src/vm/program"
)

// syncDeletedBranchProgram adds opcodes that sync a branch that was deleted at origin to the given program.
func syncDeletedBranchProgram(list *program.Program, branch undodomain.BranchInfo, parentOtherWorktree bool, args SyncBranchProgramArgs) {
	if args.BranchTypes.IsFeatureBranch(branch.LocalName) {
		syncDeletedFeatureBranchProgram(list, branch, parentOtherWorktree, args)
	} else {
		syncDeletedPerennialBranchProgram(list, branch, args)
	}
}

// syncDeletedFeatureBranchProgram syncs a feare branch whose remote has been deleted.
// The parent branch must have been fully synced before calling this function.
func syncDeletedFeatureBranchProgram(list *program.Program, branch undodomain.BranchInfo, parentOtherWorktree bool, args SyncBranchProgramArgs) {
	list.Add(&opcode.Checkout{Branch: branch.LocalName})
	pullParentBranchOfCurrentFeatureBranchOpcode(list, branch.LocalName, parentOtherWorktree, args.SyncFeatureStrategy)
	list.Add(&opcode.DeleteBranchIfEmptyAtRuntime{Branch: branch.LocalName})
}

func syncDeletedPerennialBranchProgram(list *program.Program, branch undodomain.BranchInfo, args SyncBranchProgramArgs) {
	RemoveBranchFromLineage(RemoveBranchFromLineageArgs{
		Program: list,
		Branch:  branch.LocalName,
		Parent:  args.MainBranch,
		Lineage: args.Lineage,
	})
	list.Add(&opcode.RemoveFromPerennialBranches{Branch: branch.LocalName})
	list.Add(&opcode.Checkout{Branch: args.MainBranch})
	list.Add(&opcode.DeleteLocalBranch{
		Branch: branch.LocalName,
		Force:  false,
	})
	list.Add(&opcode.QueueMessage{Message: fmt.Sprintf(messages.BranchDeleted, branch.LocalName)})
}
