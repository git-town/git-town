package opcode

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/vm/shared"
)

type RemoveLocalConfig struct {
	Key configdomain.Key // the config key to remove
	undeclaredOpcodeMethods
}

func (self *RemoveLocalConfig) Run(args shared.RunArgs) error {
	return args.Runner.GitTown.RemoveLocalConfigValue(self.Key)
}
