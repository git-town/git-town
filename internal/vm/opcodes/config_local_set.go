package opcodes

import (
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

type ConfigLocalSet struct {
	Key                     configdomain.Key
	Value                   string
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ConfigLocalSet) Run(args shared.RunArgs) error {
	return args.Config.Value.NormalConfig.GitConfig.SetConfigValue(configdomain.ConfigScopeLocal, self.Key, self.Value)
}
