package opcodes

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// CherryPick cherry-picks the given commit.
type CherryPick struct {
	SHA gitdomain.SHA
}

func (self *CherryPick) Abort() []shared.Opcode {
	return []shared.Opcode{
		&CherryPickAbort{},
	}
}

func (self *CherryPick) Continue() []shared.Opcode {
	return []shared.Opcode{
		&CherryPickContinue{},
	}
}

func (self *CherryPick) Run(args shared.RunArgs) error {
	return args.Git.CherryPick(args.Frontend, self.SHA)
}
