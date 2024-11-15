package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// RebaseBranch rebases the current branch
// against the branch with the given name.
type RebaseBranchOnto struct {
	Branch                  gitdomain.BranchName
	Onto                    gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseBranchOnto) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&RebaseAbort{},
	}
}

func (self *RebaseBranchOnto) ContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		&RebaseContinueIfNeeded{},
	}
}

func (self *RebaseBranchOnto) Run(args shared.RunArgs) error {
	return args.Git.RebaseOnto(args.Frontend, self.Branch, args.Config.Value.NormalConfig.GitVersion)
}
