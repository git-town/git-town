package opcodes

import "github.com/git-town/git-town/v21/internal/vm/shared"

// FetchUpstream brings the Git history of the local repository
// up to speed with activities that happened in the upstream remote.
type ExecuteShellCommand struct {
	Args                    []string
	Executable              string
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ExecuteShellCommand) Run(args shared.RunArgs) error {
	return args.Frontend.Run(self.Executable, self.Args...)
}
