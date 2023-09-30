package steps

import (
	"github.com/git-town/git-town/v9/src/config"
)

type SetLocalConfigStep struct {
	Key   config.Key
	Value string
	EmptyStep
}

func (step *SetLocalConfigStep) Run(args RunArgs) error {
	return args.Runner.Config.SetLocalConfigValue(step.Key, step.Value)
}
