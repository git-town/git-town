package steps

import (
	"github.com/git-town/git-town/v9/src/config"
)

type RemoveLocalConfigStep struct {
	Key config.Key // the config key to remove
	EmptyStep
}

func (step *RemoveLocalConfigStep) Run(args RunArgs) error {
	return args.Runner.Config.RemoveLocalConfigValue(step.Key)
}
