package swap

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type swapLineageParentSetsProgramArg struct {
	children    []swapBranch
	current     gitdomain.LocalBranchName
	grandParent gitdomain.LocalBranchName
	parent      gitdomain.LocalBranchName
	program     Mutable[program.Program]
}

func swapLineageParentSetsProgram(args swapLineageParentSetsProgramArg) {
	args.program.Value.Add(
		&opcodes.LineageParentSet{
			Branch: args.current,
			Parent: args.grandParent,
		},
		&opcodes.LineageParentSet{
			Branch: args.parent,
			Parent: args.current,
		},
	)
	for _, child := range args.children {
		args.program.Value.Add(
			&opcodes.LineageParentSet{
				Branch: child.name,
				Parent: args.parent,
			},
		)
	}
}
