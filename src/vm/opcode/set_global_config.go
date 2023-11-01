package opcode

import (
	"github.com/git-town/git-town/v10/src/config"
	"github.com/git-town/git-town/v10/src/vm/shared"
)

type SetGlobalConfig struct {
	Key   config.Key
	Value string
	undeclaredOpcodeMethods
}

func (self *SetGlobalConfig) Run(args shared.RunArgs) error {
	return args.Runner.Config.SetGlobalConfigValue(self.Key, self.Value)
}
