package opcodes

import (
	"github.com/git-town/git-town/v13/src/config/gitconfig"
	"github.com/git-town/git-town/v13/src/vm/shared"
)

type SetGlobalConfig struct {
	Key   gitconfig.Key
	Value string
	undeclaredOpcodeMethods
}

func (self *SetGlobalConfig) Run(args shared.RunArgs) error {
	return args.Runner.Config.GitConfig.SetGlobalConfigValue(self.Key, self.Value)
}
