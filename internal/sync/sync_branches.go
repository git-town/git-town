package sync

import (
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// BranchesProgram syncs all given branches.
func BranchesProgram(branchesToSync []configdomain.BranchToSync, args BranchProgramArgs) {
	lastAncestorSynced := None[gitdomain.LocalBranchName]()
	for _, branchToSync := range branchesToSync {
		if localName, hasLocalName := branchToSync.BranchInfo.LocalName.Get(); hasLocalName {
			LocalBranchProgram(localName, branchToSync.BranchInfo, branchToSync.FirstCommitMessage, lastAncestorSynced, BranchProgramArgs{
				BranchInfos:   args.BranchInfos,
				Config:        args.Config,
				InitialBranch: args.InitialBranch,
				Program:       args.Program,
				PushBranches:  args.PushBranches,
				Remotes:       args.Remotes,
			})
			lastAncestorSynced = Some(localName)
		}
	}
}
