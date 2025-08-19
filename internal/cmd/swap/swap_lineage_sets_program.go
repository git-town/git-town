package swap

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/opcodes"
	"github.com/git-town/git-town/v21/internal/vm/program"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type swapLineageParentSetsProgramArg struct {
	branchToSwap      gitdomain.LocalBranchName
	childBranches     []swapChildBranch
	grandParentBranch gitdomain.LocalBranchName
	parentBranch      gitdomain.LocalBranchName
	program           Mutable[program.Program]
}

func swapLineageParentSetsProgram(args swapLineageParentSetsProgramArg) {
	args.program.Value.Add(
		&opcodes.LineageParentSet{
			Branch: args.branchToSwap,
			Parent: args.grandParentBranch,
		},
		&opcodes.LineageParentSet{
			Branch: args.parentBranch,
			Parent: args.branchToSwap,
		},
	)
	for _, child := range args.childBranches {
		args.program.Value.Add(
			&opcodes.LineageParentSet{
				Branch: child.name,
				Parent: args.parentBranch,
			},
		)
	}
}
