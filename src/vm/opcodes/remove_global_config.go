package opcodes

import (
	"github.com/git-town/git-town/v12/src/config/gitconfig"
	"github.com/git-town/git-town/v12/src/vm/shared"
)

type RemoveGlobalConfig struct {
	Key gitconfig.Key // the config key to remove
	undeclaredOpcodeMethods
}

func (self *RemoveGlobalConfig) Run(args shared.RunArgs) error {
	return args.Runner.Config.GitConfig.RemoveGlobalConfigValue(self.Key)
}
