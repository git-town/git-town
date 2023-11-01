package opcode

import (
	"github.com/git-town/git-town/v10/src/config"
	"github.com/git-town/git-town/v10/src/vm/shared"
)

type RemoveGlobalConfig struct {
	Key config.Key // the config key to remove
	undeclaredOpcodeMethods
}

func (self *RemoveGlobalConfig) Run(args shared.RunArgs) error {
	return args.Runner.Config.RemoveGlobalConfigValue(self.Key)
}
