package sync

import (
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/vm/opcodes"
)

// PerennialBranchProgram adds the opcodes to sync the perennial branch with the given name.
func ObservedBranchProgram(branch gitdomain.BranchInfo, args BranchProgramArgs) {
	if branch.HasTrackingBranch() {
		args.Program.Add(&opcodes.RebaseBranch{Branch: branch.RemoteName.BranchName()})
	}
}
