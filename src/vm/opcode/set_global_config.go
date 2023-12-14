package opcode

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/vm/shared"
)

type SetGlobalConfig struct {
	Key   configdomain.Key
	Value string
	undeclaredOpcodeMethods
}

func (self *SetGlobalConfig) Run(args shared.RunArgs) error {
	return args.Runner.GitTown.SetGlobalConfigValue(self.Key, self.Value)
}
