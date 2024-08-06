package opcodes

import (
	"github.com/git-town/git-town/v14/internal/vm/shared"
	"github.com/git-town/git-town/v14/pkg/keys"
)

type SetLocalConfig struct {
	Key                     keys.Key
	Value                   string
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *SetLocalConfig) Run(args shared.RunArgs) error {
	return args.Config.GitConfig.SetLocalConfigValue(self.Key, self.Value)
}
