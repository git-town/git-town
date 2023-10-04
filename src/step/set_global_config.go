package step

import (
	"github.com/git-town/git-town/v9/src/config"
)

type SetGlobalConfig struct {
	Key   config.Key
	Value string
	Empty
}

func (step *SetGlobalConfig) Run(args RunArgs) error {
	return args.Runner.Config.SetGlobalConfigValue(step.Key, step.Value)
}
