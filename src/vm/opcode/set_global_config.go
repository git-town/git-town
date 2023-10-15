package opcode

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

type SetGlobalConfig struct {
	Key   config.Key
	Value string
	undeclaredOpcodeMethods
}

func (step *SetGlobalConfig) Run(args shared.RunArgs) error {
	return args.Runner.Config.SetGlobalConfigValue(step.Key, step.Value)
}
