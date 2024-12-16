package sync

import (
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/vm/opcodes"
	"github.com/git-town/git-town/v17/internal/vm/program"
	. "github.com/git-town/git-town/v17/pkg/prelude"
)

// PerennialBranchProgram adds the opcodes to sync the observed branch with the given name.
func ObservedBranchProgram(branch Option[gitdomain.RemoteBranchName], prog Mutable[program.Program]) {
	if remoteBranch, hasRemoteBranch := branch.Get(); hasRemoteBranch {
		prog.Value.Add(&opcodes.RebaseBranch{Branch: remoteBranch.BranchName()})
	}
}
