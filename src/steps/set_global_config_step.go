package steps

import (
	"github.com/git-town/git-town/v9/src/config"
)

type SetGlobalConfigStep struct {
	Key   config.Key
	Value string
	EmptyStep
}

func (step *SetGlobalConfigStep) Run(args RunArgs) error {
	return args.Runner.Config.SetGlobalConfigValue(step.Key, step.Value)
}
