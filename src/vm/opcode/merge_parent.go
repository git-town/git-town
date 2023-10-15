package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// MergeParent merges the branch that at runtime is the parent branch of the given branch into the given branch.
type MergeParent struct {
	CurrentBranch domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (op *MergeParent) CreateAbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&AbortMerge{},
	}
}

func (op *MergeParent) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		&ContinueMerge{},
	}
}

func (op *MergeParent) Run(args shared.RunArgs) error {
	parent := args.Lineage.Parent(op.CurrentBranch)
	if parent.IsEmpty() {
		return nil
	}
	return args.Runner.Frontend.MergeBranchNoEdit(parent.BranchName())
}
