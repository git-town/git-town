package steps

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/git"
)

type RemoveLocalConfigStep struct {
	Key           config.Key // the config key to remove
	previousValue string
	EmptyStep
}

func (step *RemoveLocalConfigStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&SetLocalConfigStep{
		Key:   step.Key,
		Value: step.previousValue,
	}}, nil
}

func (step *RemoveLocalConfigStep) Run(args RunArgs) error {
	step.previousValue = args.Runner.Config.LocalConfigValue(step.Key)
	return args.Runner.Config.RemoveLocalConfigValue(step.Key)
}
