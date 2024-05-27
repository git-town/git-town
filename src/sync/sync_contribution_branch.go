package sync

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/opcodes"
	"github.com/git-town/git-town/v14/src/vm/program"
)

// FeatureBranchProgram adds the opcodes to sync the feature branch with the given name.
func ContributionBranchProgram(prog *program.Program, branch gitdomain.BranchInfo) {
	if trackingBranch, hasTrackingBranch := branch.RemoteName.Get(); hasTrackingBranch {
		prog.Add(&opcodes.RebaseBranch{Branch: trackingBranch.BranchName()})
	}
}
