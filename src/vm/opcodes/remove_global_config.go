package opcodes

import (
	"github.com/git-town/git-town/v14/src/config/gitconfig"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

type RemoveGlobalConfig struct {
	Key gitconfig.Key // the config key to remove
	undeclaredOpcodeMethods
}

func (self *RemoveGlobalConfig) Run(args shared.RunArgs) error {
	return args.Config.GitConfig.RemoveGlobalConfigValue(self.Key)
}
