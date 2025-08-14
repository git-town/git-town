package sync

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
)

/*

New sync workflow:

- sync each stack strictly from root to leafs
- when syncing a branch:
  - sync-onto its parent to update the branch and remove the old commits
	  - the old commits are:
		  - what the branch had at the end of the previous sync, if it's different from what it had at the beginning of the sync
			- what the branch had before the sync ran, but only if it's different from its current SHA
			- if both are available: which one to use?
		- if no commits to remove, rebase against the parent normally
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
}
