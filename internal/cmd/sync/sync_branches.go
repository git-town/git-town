package sync

import (
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/vm/opcodes"
)

// BranchesProgram syncs all given branches.
func BranchesProgram(branchesToSync []configdomain.BranchToSync, args BranchProgramArgs) {
	for _, branchToSync := range branchesToSync {
		if localBranchName, hasLocalBranch := branchToSync.BranchInfo.LocalName.Get(); hasLocalBranch {
			BranchProgram(localBranchName, branchToSync.BranchInfo, branchToSync.FirstCommitMessage, args)
		}
	}
	if parentToDelete, hasParentToDelete := args.ParentToDelete.Get(); hasParentToDelete {
		args.Program.Value.Add(&opcodes.BranchLocalDelete{
			Branch: parentToDelete,
		})
	}
}
