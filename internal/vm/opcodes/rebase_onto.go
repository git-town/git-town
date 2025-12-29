package opcodes

import (
	"time"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// rebases the current branch against the target branch while executing "git town swap", while moving the target branch onto the Onto branch.
type RebaseOnto struct {
	BranchToRebaseOnto gitdomain.BranchName
	CommitsToRemove    gitdomain.Location
}

func (self *RebaseOnto) Abort() []shared.Opcode {
	return []shared.Opcode{
		&RebaseAbort{},
	}
}

func (self *RebaseOnto) Continue() []shared.Opcode {
	return []shared.Opcode{
		&RebaseContinueIfNeeded{},
	}
}

func (self *RebaseOnto) Run(args shared.RunArgs) error {
	if args.Config.Value.NormalConfig.TestHome.IsSome() {
		// Fix for https://github.com/git-town/git-town/issues/4942.
		// Waiting here in end-to-end tests to ensure new timestamps for the rebased commits,
		// which avoids flaky end-to-end tests.
		time.Sleep(1 * time.Second)
	}
	return args.Git.RebaseOnto(args.Frontend, self.BranchToRebaseOnto.Location(), self.CommitsToRemove)
}
