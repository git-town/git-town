package opcodes

import (
	"github.com/git-town/git-town/v12/src/config/gitconfig"
	"github.com/git-town/git-town/v12/src/vm/shared"
)

type RemoveLocalConfig struct {
	Key gitconfig.Key // the config key to remove
	undeclaredOpcodeMethods
}

func (self *RemoveLocalConfig) Run(args shared.RunArgs) error {
	return args.Runner.Config.GitConfig.RemoveLocalConfigValue(self.Key)
}
