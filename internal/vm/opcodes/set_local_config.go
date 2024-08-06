package opcodes

import (
	"github.com/git-town/git-town/v14/internal/config/configdomain"
	"github.com/git-town/git-town/v14/internal/vm/shared"
)

type SetLocalConfig struct {
	Key                     configdomain.Key
	Value                   string
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *SetLocalConfig) Run(args shared.RunArgs) error {
	return args.Config.GitConfig.SetLocalConfigValue(self.Key, self.Value)
}
