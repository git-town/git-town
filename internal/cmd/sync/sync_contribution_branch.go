package sync

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// ContributionBranchProgram adds the opcodes to sync the feature branch with the given name.
func ContributionBranchProgram(prog Mutable[program.Program], branchInfo gitdomain.BranchInfo) {
	if trackingBranch, hasTrackingBranch := branchInfo.RemoteName.Get(); hasTrackingBranch {
		if branchInfo.SyncStatus != gitdomain.SyncStatusUpToDate {
			prog.Value.Add(&opcodes.RebaseBranch{Branch: trackingBranch.BranchName()})
		}
	}
}
