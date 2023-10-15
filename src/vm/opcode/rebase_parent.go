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

func (step *RebaseParent) CreateAbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&AbortRebase{},
	}
}

func (step *RebaseParent) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		&ContinueRebase{},
	}
}

func (step *RebaseParent) Run(args shared.RunArgs) error {
	parent := args.Lineage.Parent(step.CurrentBranch)
	if parent.IsEmpty() {
		return nil
	}
	return args.Runner.Frontend.Rebase(parent.BranchName())
}
