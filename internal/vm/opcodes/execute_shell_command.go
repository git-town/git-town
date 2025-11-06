package opcodes

import "github.com/git-town/git-town/v22/internal/vm/shared"

// ExecuteShellCommand executes a shell command.
type ExecuteShellCommand struct {
	Args       []string
	Executable string
}

func (self *ExecuteShellCommand) Run(args shared.RunArgs) error {
	return args.Frontend.Run(self.Executable, self.Args...)
}
