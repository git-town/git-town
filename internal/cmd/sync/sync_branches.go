package sync

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/vm/opcodes"
)

// BranchesProgram syncs all given branches.
func BranchesProgram(branchesToSync []configdomain.BranchToSync, args BranchProgramArgs) {
	for _, branchToSync := range branchesToSync {
		if localBranchName, hasLocalBranch := branchToSync.BranchInfo.LocalName.Get(); hasLocalBranch {
			BranchProgram(localBranchName, branchToSync.BranchInfo, branchToSync.FirstCommitMessage, args)
			fmt.Println("333333333333 parentToDelete", args.BranchesToDelete)
		}
	}
	for _, branchToDelete := range args.BranchesToDelete.Value.Values() {
		args.Program.Value.Add(
			&opcodes.BranchLocalDelete{Branch: branchToDelete},
			&opcodes.LineageBranchRemove{Branch: branchToDelete},
		)
	}
}
