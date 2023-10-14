package step

import (
	"github.com/git-town/git-town/v9/src/config"
)

type RemoveLocalConfig struct {
	Key config.Key // the config key to remove
	Empty
}

func (step *RemoveLocalConfig) Run(args RunArgs) error {
	return args.Runner.Config.RemoveLocalConfigValue(step.Key)
}
