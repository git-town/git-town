package opcodes

import (
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

// Checkout checks out the given existing branch.
type CherryPickAbort struct {
}

func (self *CherryPickAbort) Run(args shared.RunArgs) error {
	return args.Git.CherryPickAbort(args.Frontend)
}
