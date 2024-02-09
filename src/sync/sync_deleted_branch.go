package sync

import (
	"fmt"

	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/git-town/git-town/v12/src/vm/opcode"
	"github.com/git-town/git-town/v12/src/vm/program"
)

// syncDeletedBranchProgram adds opcodes that sync a branch that was deleted at origin to the given program.
func syncDeletedBranchProgram(list *program.Program, branch gitdomain.BranchInfo, parentOtherWorktree bool, args BranchProgramArgs) {
	if args.Config.IsFeatureBranch(branch.LocalName) {
		syncDeletedFeatureBranchProgram(list, branch, parentOtherWorktree, args)
	} else {
		syncDeletedPerennialBranchProgram(list, branch, args)
	}
}

// syncDeletedFeatureBranchProgram syncs a feare branch whose remote has been deleted.
// The parent branch must have been fully synced before calling this function.
func syncDeletedFeatureBranchProgram(list *program.Program, branch gitdomain.BranchInfo, parentOtherWorktree bool, args BranchProgramArgs) {
	list.Add(&opcode.Checkout{Branch: branch.LocalName})
	pullParentBranchOfCurrentFeatureBranchOpcode(list, branch.LocalName, parentOtherWorktree, args.Config.SyncFeatureStrategy)
	list.Add(&opcode.DeleteBranchIfEmptyAtRuntime{Branch: branch.LocalName})
}

func syncDeletedPerennialBranchProgram(list *program.Program, branch gitdomain.BranchInfo, args BranchProgramArgs) {
	RemoveBranchFromLineage(RemoveBranchFromLineageArgs{
		Branch:  branch.LocalName,
		Lineage: args.Config.Lineage,
		Parent:  args.Config.MainBranch,
		Program: list,
	})
	list.Add(&opcode.RemoveFromPerennialBranches{Branch: branch.LocalName})
	list.Add(&opcode.Checkout{Branch: args.Config.MainBranch})
	list.Add(&opcode.DeleteLocalBranch{
		Branch: branch.LocalName,
		Force:  false,
	})
	list.Add(&opcode.QueueMessage{Message: fmt.Sprintf(messages.BranchDeleted, branch.LocalName)})
}
