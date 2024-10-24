package opcodes

import (
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

type ConfigLocalRemove struct {
	Key                     configdomain.Key // the config key to remove
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ConfigLocalRemove) Run(args shared.RunArgs) error {
	return args.Config.Value.NormalConfig.GitConfig.RemoveLocalConfigValue(self.Key)
}
