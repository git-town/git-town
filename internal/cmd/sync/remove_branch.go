package sync

import (
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/vm/opcodes"
	"github.com/git-town/git-town/v17/internal/vm/program"
	. "github.com/git-town/git-town/v17/pkg/prelude"
)

func RemoveBranchConfiguration(args RemoveBranchConfigurationArgs) {
	args.Program.Value.Add(&opcodes.BranchesContributionRemove{Branch: args.Branch})
	args.Program.Value.Add(&opcodes.BranchesObservedRemove{Branch: args.Branch})
	args.Program.Value.Add(&opcodes.BranchesParkedRemove{Branch: args.Branch})
	args.Program.Value.Add(&opcodes.BranchesPerennialRemove{Branch: args.Branch})
	args.Program.Value.Add(&opcodes.BranchesPrototypeRemove{Branch: args.Branch})
	childBranches := args.Lineage.Children(args.Branch)
	for _, child := range childBranches {
		args.Program.Value.Add(&opcodes.LineageParentSetToGrandParent{Branch: child})
	}
	args.Program.Value.Add(&opcodes.LineageParentRemove{Branch: args.Branch})
}

type RemoveBranchConfigurationArgs struct {
	Branch  gitdomain.LocalBranchName
	Lineage configdomain.Lineage
	Program Mutable[program.Program]
}
