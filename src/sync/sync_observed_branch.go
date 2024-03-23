package sync

import (
	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/git-town/git-town/v13/src/vm/opcodes"
	"github.com/git-town/git-town/v13/src/vm/program"
)

// PerennialBranchProgram adds the opcodes to sync the observed branch with the given name.
func ObservedBranchProgram(branch gitdomain.BranchInfo, prog *program.Program) {
	if branch.HasTrackingBranch() {
		prog.Add(&opcodes.RebaseBranch{Branch: branch.RemoteName.BranchName()})
	}
}
