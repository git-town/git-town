package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// Merge merges the branch with the given name into the current branch.
type Merge struct {
	Branch domain.BranchName
	undeclaredOpcodeMethods
}

func (op *Merge) CreateAbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&AbortMerge{},
	}
}

func (op *Merge) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		&ContinueMerge{},
	}
}

func (op *Merge) Run(args shared.RunArgs) error {
	return args.Runner.Frontend.MergeBranchNoEdit(op.Branch)
}
