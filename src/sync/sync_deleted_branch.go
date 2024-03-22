package sync

import (
	"fmt"

	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/git-town/git-town/v12/src/vm/opcodes"
	"github.com/git-town/git-town/v12/src/vm/program"
)

// syncDeletedBranchProgram adds opcodes that sync a branch that was deleted at origin to the given program.
func syncDeletedBranchProgram(list *program.Program, branch gitdomain.BranchInfo, parentOtherWorktree bool, args BranchProgramArgs) {
	switch args.Config.BranchType(branch.LocalName) {
	case configdomain.BranchTypeFeatureBranch:
		syncDeletedFeatureBranchProgram(list, branch, parentOtherWorktree, args)
	case configdomain.BranchTypePerennialBranch, configdomain.BranchTypeMainBranch:
		syncDeletedPerennialBranchProgram(list, branch, args)
	case configdomain.BranchTypeObservedBranch, configdomain.BranchTypeContributionBranch, configdomain.BranchTypeParkedBranch:
		syncDeletedObservedBranchProgram(list, branch, args)
	}
}

// syncDeletedFeatureBranchProgram syncs a feare branch whose remote has been deleted.
// The parent branch must have been fully synced before calling this function.
func syncDeletedFeatureBranchProgram(list *program.Program, branch gitdomain.BranchInfo, parentOtherWorktree bool, args BranchProgramArgs) {
	list.Add(&opcodes.Checkout{Branch: branch.LocalName})
	pullParentBranchOfCurrentFeatureBranchOpcode(featureBranchArgs{
		branch:              branch,
		offline:             args.Config.Offline,
		parentOtherWorktree: parentOtherWorktree,
		program:             list,
		syncStrategy:        args.Config.SyncFeatureStrategy,
	})
	list.Add(&opcodes.DeleteBranchIfEmptyAtRuntime{Branch: branch.LocalName})
}

func syncDeletedObservedBranchProgram(list *program.Program, branch gitdomain.BranchInfo, args BranchProgramArgs) {
	RemoveBranchFromLineage(RemoveBranchFromLineageArgs{
		Branch:  branch.LocalName,
		Lineage: args.Config.Lineage,
		Parent:  args.Config.MainBranch,
		Program: list,
	})
	list.Add(&opcodes.RemoveFromObservedBranches{Branch: branch.LocalName})
	list.Add(&opcodes.Checkout{Branch: args.Config.MainBranch})
	list.Add(&opcodes.DeleteLocalBranch{Branch: branch.LocalName})
	list.Add(&opcodes.QueueMessage{Message: fmt.Sprintf(messages.BranchDeleted, branch.LocalName)})
}

func syncDeletedPerennialBranchProgram(list *program.Program, branch gitdomain.BranchInfo, args BranchProgramArgs) {
	RemoveBranchFromLineage(RemoveBranchFromLineageArgs{
		Branch:  branch.LocalName,
		Lineage: args.Config.Lineage,
		Parent:  args.Config.MainBranch,
		Program: list,
	})
	list.Add(&opcodes.RemoveFromPerennialBranches{Branch: branch.LocalName})
	list.Add(&opcodes.Checkout{Branch: args.Config.MainBranch})
	list.Add(&opcodes.DeleteLocalBranch{Branch: branch.LocalName})
	list.Add(&opcodes.QueueMessage{Message: fmt.Sprintf(messages.BranchDeleted, branch.LocalName)})
}
