package opcodes

import (
	"github.com/git-town/git-town/v14/internal/vm/shared"
	"github.com/git-town/git-town/v14/pkg/keys"
)

type SetGlobalConfig struct {
	Key                     keys.Key
	Value                   string
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *SetGlobalConfig) Run(args shared.RunArgs) error {
	return args.Config.GitConfig.SetGlobalConfigValue(self.Key, self.Value)
}
