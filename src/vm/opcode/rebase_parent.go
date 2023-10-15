package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// RebaseParent rebases the given branch against the branch that is its parent at runtime.
type RebaseParent struct {
	CurrentBranch domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (op *RebaseParent) CreateAbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&AbortRebase{},
	}
}

func (op *RebaseParent) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		&ContinueRebase{},
	}
}

func (op *RebaseParent) Run(args shared.RunArgs) error {
	parent := args.Lineage.Parent(op.CurrentBranch)
	if parent.IsEmpty() {
		return nil
	}
	return args.Runner.Frontend.Rebase(parent.BranchName())
}
