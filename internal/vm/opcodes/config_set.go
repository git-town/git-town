package opcodes

import (
	"github.com/git-town/git-town/v18/internal/config/configdomain"
	"github.com/git-town/git-town/v18/internal/vm/shared"
)

type ConfigSet struct {
	Key                     configdomain.Key
	Scope                   configdomain.ConfigScope
	Value                   string
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ConfigSet) Run(args shared.RunArgs) error {
	return args.Config.Value.NormalConfig.GitConfigAccess.SetConfigValue(self.Scope, self.Key, self.Value)
}
