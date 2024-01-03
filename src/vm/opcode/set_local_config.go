package opcode

import (
	"github.com/git-town/git-town/v11/src/config/gitconfig"
	"github.com/git-town/git-town/v11/src/vm/shared"
)

type SetLocalConfig struct {
	Key   gitconfig.Key
	Value string
	undeclaredOpcodeMethods
}

func (self *SetLocalConfig) Run(args shared.RunArgs) error {
	return args.Runner.Config.Access.SetLocalConfigValue(self.Key, self.Value)
}
