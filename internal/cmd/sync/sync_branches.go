package sync

import "github.com/git-town/git-town/v16/internal/config/configdomain"

// BranchesProgram syncs all given branches.
func BranchesProgram(branchesToSync []configdomain.BranchToSync, args BranchProgramArgs) {
	for _, branchToSync := range branchesToSync {
		if localBranchName, hasLocalBranch := branchToSync.BranchInfo.LocalName.Get(); hasLocalBranch {
			BranchProgram(localBranchName, branchToSync.BranchInfo, branchToSync.FirstCommitMessage, args)
		}
	}
}
