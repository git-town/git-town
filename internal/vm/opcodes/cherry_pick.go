package opcodes

import (
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/vm/shared"
)

// Checkout checks out the given existing branch.
type CherryPick struct {
	SHA                     gitdomain.SHA
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CherryPick) Run(args shared.RunArgs) error {
	return args.Git.CherryPick(args.Frontend, self.SHA)
}
