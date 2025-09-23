package sync

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// ObservedBranchProgram adds the opcodes to sync the observed branch with the given name.
func ObservedBranchProgram(branchInfo gitdomain.BranchInfo, prog Mutable[program.Program]) {
	if remoteBranch, hasRemoteBranch := branchInfo.RemoteName.Get(); hasRemoteBranch {
		if branchInfo.SyncStatus != gitdomain.SyncStatusUpToDate {
			prog.Value.Add(&opcodes.RebaseBranch{Branch: remoteBranch.BranchName()})
		}
	}
}
