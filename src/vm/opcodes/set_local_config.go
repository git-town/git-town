package opcodes

import (
	"github.com/git-town/git-town/v14/src/config/gitconfig"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

type SetLocalConfig struct {
	Key   gitconfig.Key
	Value string
	undeclaredOpcodeMethods
}

func (self *SetLocalConfig) Run(args shared.RunArgs) error {
	return args.Config.GitConfig.SetLocalConfigValue(self.Key, self.Value)
}
