package opcode

import (
	"github.com/git-town/git-town/v10/src/config"
	"github.com/git-town/git-town/v10/src/vm/shared"
)

type RemoveLocalConfig struct {
	Key config.Key // the config key to remove
	undeclaredOpcodeMethods
}

func (self *RemoveLocalConfig) Run(args shared.RunArgs) error {
	return args.Runner.Config.RemoveLocalConfigValue(self.Key)
}
