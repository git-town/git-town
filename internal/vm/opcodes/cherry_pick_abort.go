package opcodes

import (
	"github.com/git-town/git-town/v20/internal/vm/shared"
)

// Checkout checks out the given existing branch.
type CherryPickAbort struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CherryPickAbort) Run(args shared.RunArgs) error {
	return args.Git.CherryPickAbort(args.Frontend)
}
