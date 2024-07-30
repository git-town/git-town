package opcodes

import (
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

type RemoveLocalConfig struct {
	Key                     configdomain.Key // the config key to remove
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RemoveLocalConfig) Run(args shared.RunArgs) error {
	return args.Config.GitConfig.RemoveLocalConfigValue(self.Key)
}
