package opcode

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/vm/shared"
)

type SetLocalConfig struct {
	Key   configdomain.Key
	Value string
	undeclaredOpcodeMethods
}

func (self *SetLocalConfig) Run(args shared.RunArgs) error {
	return args.Runner.GitTown.SetLocalConfigValue(self.Key, self.Value)
}
