package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// RebaseBranch rebases the current branch
// against the branch with the given name.
type RebaseBranch struct {
	Branch domain.BranchName
	undeclaredOpcodeMethods
}

func (step *RebaseBranch) CreateAbortProgram() []Opcode {
	return []Opcode{&AbortRebase{}}
}

func (step *RebaseBranch) CreateContinueProgram() []Opcode {
	return []Opcode{&ContinueRebase{}}
}

func (step *RebaseBranch) Run(args RunArgs) error {
	return args.Runner.Frontend.Rebase(step.Branch)
}
