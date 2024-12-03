package sync

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/opcodes"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// BranchesProgram syncs all given branches.
func BranchesProgram(branchesToSync []configdomain.BranchToSync, args Mutable[BranchProgramArgs]) {
	for _, branchToSync := range branchesToSync {
		if localBranchName, hasLocalBranch := branchToSync.BranchInfo.LocalName.Get(); hasLocalBranch {
			BranchProgram(localBranchName, branchToSync.BranchInfo, branchToSync.FirstCommitMessage, args)
			fmt.Println("333333333333333333333333333333333333333333", args.Value.ParentToDelete)
		}
	}
	if parentToDelete, hasParentToDelete := args.Value.ParentToDelete.Get(); hasParentToDelete {
		args.Value.Program.Value.Add(
			&opcodes.BranchLocalDelete{
				Branch: parentToDelete,
			},
			&opcodes.LineageBranchRemove{
				Branch: parentToDelete,
			},
		)
		args.Value.ParentToDelete = None[gitdomain.LocalBranchName]()
	}
}
