package opcodes

import (
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

// Checkout checks out the given existing branch.
type CherryPickContinue struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CherryPickContinue) Run(args shared.RunArgs) error {
	return args.Git.CherryPickContinue(args.Frontend)
}
