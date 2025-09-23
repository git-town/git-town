package opcodes

import "github.com/git-town/git-town/v22/internal/vm/shared"

// RebaseContinue finishes an ongoing rebase operation
// assuming all conflicts have been resolved by the user.
type RebaseContinue struct{}

func (self *RebaseContinue) Abort() []shared.Opcode {
	return []shared.Opcode{
		&RebaseAbort{},
	}
}

func (self *RebaseContinue) Run(args shared.RunArgs) error {
	return args.Git.ContinueRebase(args.Frontend)
}
