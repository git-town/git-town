package opcodes

import (
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/vm/shared"
)

type SetGlobalConfig struct {
	Key                     configdomain.Key
	Value                   string
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *SetGlobalConfig) Run(args shared.RunArgs) error {
	return args.Config.GitConfig.SetGlobalConfigValue(self.Key, self.Value)
}
