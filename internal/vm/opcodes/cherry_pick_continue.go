package opcodes

import (
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// CherryPickContinue continues a suspended cherry-pick operation.
type CherryPickContinue struct{}

func (self *CherryPickContinue) Run(args shared.RunArgs) error {
	return args.Git.CherryPickContinue(args.Frontend)
}
