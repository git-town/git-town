package opcodes

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/subshell"
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

// RebaseBranch rebases the current branch
// against the branch with the given name.
type RebaseBranch struct {
	Branch                  gitdomain.BranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseBranch) Abort() []shared.Opcode {
	return []shared.Opcode{
		&RebaseAbort{},
	}
}

func (self *RebaseBranch) Continue() []shared.Opcode {
	return []shared.Opcode{
		&RebaseContinueIfNeeded{},
	}
}

func (self *RebaseBranch) Run(args shared.RunArgs) error {
	// Fix for https://github.com/git-town/git-town/issues/4942.
	// Waiting here in end-to-end tests to ensure new timestamps for the rebased commits,
	// which avoids flaky end-to-end tests.
	if subshell.IsInTest() {
		args.Frontend.Run("sleep", "1")
		// time.Sleep(1 * time.Second)
	}
	return args.Git.Rebase(args.Frontend, self.Branch)
}
