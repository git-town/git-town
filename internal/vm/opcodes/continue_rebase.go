package opcodes

import "github.com/git-town/git-town/v16/internal/vm/shared"

// ContinueRebase finishes an ongoing rebase operation
// assuming all conflicts have been resolved by the user.
type ContinueRebase struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ContinueRebase) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&AbortRebase{},
	}
}

func (self *ContinueRebase) Run(args shared.RunArgs) error {
	return args.Git.ContinueRebase(args.Frontend)
}
