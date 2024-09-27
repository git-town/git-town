package sync

import "github.com/git-town/git-town/v16/internal/config/configdomain"

// BranchesProgram syncs all given branches.
func BranchesProgram(branchesToSync []configdomain.BranchToSync, args BranchProgramArgs) {
	for _, branchToSync := range branchesToSync {
		BranchProgram(branchToSync.BranchInfo, branchToSync.FirstCommitMessage, BranchProgramArgs{
			BranchInfos:   args.BranchInfos,
			Config:        args.Config,
			InitialBranch: args.InitialBranch,
			Program:       args.Program,
			PushBranches:  args.PushBranches,
			Remotes:       args.Remotes,
		})
	}
}
