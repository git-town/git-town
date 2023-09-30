package steps

import (
	"github.com/git-town/git-town/v9/src/config"
)

type RemoveGlobalConfigStep struct {
	Key config.Key // the config key to remove
	EmptyStep
}

func (step *RemoveGlobalConfigStep) Run(args RunArgs) error {
	return args.Runner.Config.RemoveGlobalConfigValue(step.Key)
}
