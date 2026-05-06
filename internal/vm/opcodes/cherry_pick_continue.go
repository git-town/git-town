package opcodes

import (
	"github.com/git-town/git-town/v23/internal/vm/shared"
)

// CherryPickContinue continues a suspended cherry-pick operation.
type CherryPickContinue struct{}

func (self *CherryPickContinue) Run(args shared.RunArgs) error {
	return args.Git.CherryPickContinue(args.Frontend)
}
