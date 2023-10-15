package opcode

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

type SetLocalConfig struct {
	Key   config.Key
	Value string
	undeclaredOpcodeMethods
}

func (op *SetLocalConfig) Run(args shared.RunArgs) error {
	return args.Runner.Config.SetLocalConfigValue(op.Key, op.Value)
}
