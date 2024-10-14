package sync

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/opcodes"
	"github.com/git-town/git-town/v16/internal/vm/program"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// FeatureBranchProgram adds the opcodes to sync the feature branch with the given name.
func ContributionBranchProgram(prog Mutable[program.Program], branch gitdomain.BranchInfo) {
	if trackingBranch, hasTrackingBranch := branch.RemoteName.Get(); hasTrackingBranch {
		prog.Value.Add(&opcodes.RebaseBranch{Branch: trackingBranch.BranchName()})
	}
}
