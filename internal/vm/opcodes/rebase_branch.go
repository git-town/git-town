package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// RebaseBranch rebases the current branch
// against the branch with the given name.
type RebaseBranch struct {
	Branch                  gitdomain.BranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseBranch) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&RebaseAbort{},
	}
}

func (self *RebaseBranch) ContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		&ContinueRebaseIfNeeded{},
	}
}

func (self *RebaseBranch) Run(args shared.RunArgs) error {
	return args.Git.Rebase(args.Frontend, self.Branch, args.Config.GitVersion)
}
