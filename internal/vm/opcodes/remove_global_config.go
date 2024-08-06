package opcodes

import (
	"github.com/git-town/git-town/v14/internal/vm/shared"
	"github.com/git-town/git-town/v14/pkg/keys"
)

type RemoveGlobalConfig struct {
	Key                     keys.Key // the config key to remove
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RemoveGlobalConfig) Run(args shared.RunArgs) error {
	return args.Config.GitConfig.RemoveGlobalConfigValue(self.Key)
}
