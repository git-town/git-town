package sync

import (
	"github.com/git-town/git-town/v14/internal/git/gitdomain"
	. "github.com/git-town/git-town/v14/internal/gohacks/prelude"
	"github.com/git-town/git-town/v14/internal/vm/opcodes"
	"github.com/git-town/git-town/v14/internal/vm/program"
)

// FeatureBranchProgram adds the opcodes to sync the feature branch with the given name.
func ContributionBranchProgram(prog Mutable[program.Program], branch gitdomain.BranchInfo) {
	if trackingBranch, hasTrackingBranch := branch.RemoteName.Get(); hasTrackingBranch {
		prog.Value.Add(&opcodes.RebaseBranch{Branch: trackingBranch.BranchName()})
	}
}
