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

func (step *Merge) CreateAbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&AbortMerge{},
	}
}

func (step *Merge) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		&ContinueMerge{},
	}
}

func (step *Merge) Run(args shared.RunArgs) error {
	return args.Runner.Frontend.MergeBranchNoEdit(step.Branch)
}
