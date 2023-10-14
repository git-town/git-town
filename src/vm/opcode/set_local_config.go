package opcode

import (
	"github.com/git-town/git-town/v9/src/config"
)

type SetLocalConfig struct {
	Key   config.Key
	Value string
	Empty
}

func (step *SetLocalConfig) Run(args RunArgs) error {
	return args.Runner.Config.SetLocalConfigValue(step.Key, step.Value)
}
