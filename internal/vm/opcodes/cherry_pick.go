package opcodes

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

// Checkout checks out the given existing branch.
type CherryPick struct {
	SHA                     gitdomain.SHA
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CherryPick) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&CherryPickAbort{},
	}
}

func (self *CherryPick) ContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		&CherryPickContinue{},
	}
}

func (self *CherryPick) Run(args shared.RunArgs) error {
	return args.Git.CherryPick(args.Frontend, self.SHA)
}
