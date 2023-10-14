package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// Merge merges the branch with the given name into the current branch.
type Merge struct {
	Branch domain.BranchName
	BaseOpcode
}

func (step *Merge) CreateAbortProgram() []Opcode {
	return []Opcode{
		&AbortMerge{},
	}
}

func (step *Merge) CreateContinueProgram() []Opcode {
	return []Opcode{
		&ContinueMerge{},
	}
}

func (step *Merge) Run(args RunArgs) error {
	return args.Runner.Frontend.MergeBranchNoEdit(step.Branch)
}
