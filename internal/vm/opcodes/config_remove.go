package opcodes

import (
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

type ConfigRemove struct {
	Key                     configdomain.Key // the config key to remove
	Scope                   configdomain.ConfigScope
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ConfigRemove) Run(args shared.RunArgs) error {
	return args.Config.Value.NormalConfig.GitConfigAccess.RemoveConfigValue(self.Scope, self.Key)
}
