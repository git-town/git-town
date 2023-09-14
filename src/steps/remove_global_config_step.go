package steps

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/git"
)

type RemoveGlobalConfigStep struct {
	Key           config.Key // the config key to remove
	previousValue string
	EmptyStep
}

func (step *RemoveGlobalConfigStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&SetGlobalConfigStep{
		Key:   step.Key,
		Value: step.previousValue,
	}}, nil
}

func (step *RemoveGlobalConfigStep) Run(args RunArgs) error {
	step.previousValue = args.Runner.Config.GlobalConfigValue(step.Key)
	return args.Runner.Config.RemoveGlobalConfigValue(step.Key)
}
