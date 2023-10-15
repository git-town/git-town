package opcode

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

type RemoveGlobalConfig struct {
	Key config.Key // the config key to remove
	undeclaredOpcodeMethods
}

func (op *RemoveGlobalConfig) Run(args shared.RunArgs) error {
	return args.Runner.Config.RemoveGlobalConfigValue(op.Key)
}
