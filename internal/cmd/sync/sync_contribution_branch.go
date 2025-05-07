package sync

import (
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/vm/opcodes"
	"github.com/git-town/git-town/v20/internal/vm/program"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

// ContributionBranchProgram adds the opcodes to sync the feature branch with the given name.
func ContributionBranchProgram(prog Mutable[program.Program], branchInfo gitdomain.BranchInfo) {
	if trackingBranch, hasTrackingBranch := branchInfo.RemoteName.Get(); hasTrackingBranch {
		if branchInfo.SyncStatus != gitdomain.SyncStatusUpToDate {
			prog.Value.Add(&opcodes.RebaseBranch{Branch: trackingBranch.BranchName()})
		}
	}
}
