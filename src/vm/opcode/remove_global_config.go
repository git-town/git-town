package opcode

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/vm/shared"
)

type RemoveGlobalConfig struct {
	Key configdomain.Key // the config key to remove
	undeclaredOpcodeMethods
}

func (self *RemoveGlobalConfig) Run(args shared.RunArgs) error {
	return args.Runner.GitTown.RemoveGlobalConfigValue(self.Key)
}
