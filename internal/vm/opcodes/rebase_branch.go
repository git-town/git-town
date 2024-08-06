package opcodes

import (
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/vm/shared"
)

// RebaseBranch rebases the current branch
// against the branch with the given name.
type RebaseBranch struct {
	Branch                  gitdomain.BranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseBranch) CreateAbortProgram() []shared.Opcode {
	return []shared.Opcode{&AbortRebase{}}
}

func (self *RebaseBranch) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		&ContinueRebase{},
	}
}

func (self *RebaseBranch) Run(args shared.RunArgs) error {
	return args.Git.Rebase(args.Frontend, self.Branch)
}
