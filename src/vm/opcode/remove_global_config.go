package opcode

import (
	"github.com/git-town/git-town/v9/src/config"
)

type RemoveGlobalConfig struct {
	Key config.Key // the config key to remove
	undeclaredOpcodeMethods
}

func (step *RemoveGlobalConfig) Run(args RunArgs) error {
	return args.Runner.Config.RemoveGlobalConfigValue(step.Key)
}
