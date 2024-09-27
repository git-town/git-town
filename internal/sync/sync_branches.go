package sync

import "github.com/git-town/git-town/v16/internal/config/configdomain"

// BranchesProgram syncs all given branches.
func BranchesProgram(args BranchesProgramArgs) {
	for _, branchToSync := range args.BranchesToSync {
		BranchProgram(branchToSync.BranchInfo, branchToSync.FirstCommitMessage, BranchProgramArgs{
			BranchInfos:   args.BranchInfos,
			Config:        args.Config,
			InitialBranch: args.InitialBranch,
			Remotes:       args.Remotes,
			Program:       args.Program,
			PushBranches:  args.PushBranches,
		})
	}
}

type BranchesProgramArgs struct {
	BranchProgramArgs
	BranchesToSync []configdomain.BranchToSync
}
