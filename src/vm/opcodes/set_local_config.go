package opcodes

import (
	"github.com/git-town/git-town/v13/src/config/gitconfig"
	"github.com/git-town/git-town/v13/src/vm/shared"
)

type SetLocalConfig struct {
	Key   gitconfig.Key
	Value string
	undeclaredOpcodeMethods
}

func (self *SetLocalConfig) Run(args shared.RunArgs) error {
	return args.Runner.Config.GitConfig.SetLocalConfigValue(self.Key, self.Value)
}
