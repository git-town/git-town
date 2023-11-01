package opcode

import (
	"github.com/git-town/git-town/v10/src/config"
	"github.com/git-town/git-town/v10/src/vm/shared"
)

type SetLocalConfig struct {
	Key   config.Key
	Value string
	undeclaredOpcodeMethods
}

func (self *SetLocalConfig) Run(args shared.RunArgs) error {
	return args.Runner.Config.SetLocalConfigValue(self.Key, self.Value)
}
