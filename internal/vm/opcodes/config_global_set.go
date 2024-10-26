package opcodes

import (
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

type ConfigGlobalSet struct {
	Key                     configdomain.Key
	Value                   string
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ConfigGlobalSet) Run(args shared.RunArgs) error {
	return args.Config.Value.NormalConfig.GitConfig.SetConfigValue(configdomain.ConfigScopeGlobal, self.Key, self.Value)
}
