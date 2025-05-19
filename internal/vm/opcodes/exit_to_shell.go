package opcodes

import (
	"github.com/git-town/git-town/v20/internal/vm/shared"
)

// FetchUpstream brings the Git history of the local repository
// up to speed with activities that happened in the upstream remote.
type ExitToShell struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ExitToShell) Run(args shared.RunArgs) error {
	return ExitToShellSignal{}
}

// ExitToShellSignal is a sentinel error that signals that no error happened
// and Git Town should simply exit to the shell without an error code,
// allowing resume via "git town continue".
type ExitToShellSignal struct{}

func (self ExitToShellSignal) Error() string {
	return ""
}
