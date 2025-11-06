package opcodes

import (
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// CherryPickAbort aborts a suspended cherry-pick operation.
type CherryPickAbort struct{}

func (self *CherryPickAbort) Run(args shared.RunArgs) error {
	return args.Git.CherryPickAbort(args.Frontend)
}
