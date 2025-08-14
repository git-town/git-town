package sync

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/vm/opcodes"
)

/*

New sync workflow:

- sync each stack strictly from root to leafs
- when syncing a branch:
  - sync-onto its parent to update the branch and remove the old commits
	  - the old commits are those that the branch had at the end of the previous sync
		- if no previous commits are known, rebase against the parent normally
	- force-push-with-lease to the tracking branch
	  - if that fails, rebase against the tracking branch and then force-push-with-lease again
- when syncing a branch whose remote is gone:
  - sync-onto its parent
	- if the branch is now empty: delete it

*/

// BranchesProgram syncs all given branches.
func BranchesProgram(branchesToSync configdomain.BranchesToSync, args BranchProgramArgs) {
	for _, branchToSync := range branchesToSync {
		if localBranchName, hasLocalBranch := branchToSync.BranchInfo.LocalName.Get(); hasLocalBranch {
			BranchProgram(localBranchName, branchToSync.BranchInfo, branchToSync.FirstCommitMessage, args)
		}
	}
	for _, branchToDelete := range args.BranchesToDelete.Value.Values() {
		args.Program.Value.Add(
			&opcodes.BranchLocalDelete{Branch: branchToDelete},
			&opcodes.LineageBranchRemove{Branch: branchToDelete},
		)
	}
}
