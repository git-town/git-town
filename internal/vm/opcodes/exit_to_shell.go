package opcodes

import (
	"github.com/git-town/git-town/v20/internal/vm/shared"
)

// FetchUpstream brings the Git history of the local repository
// up to speed with activities that happened in the upstream remote.
type ExitToShell struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ExitToShell) Run(_ shared.RunArgs) error {
	return shared.ErrExitToShell
}
