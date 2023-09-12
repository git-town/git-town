package steps

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
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

func (step *RemoveLocalConfigStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	step.previousValue = run.Config.LocalConfigValue(step.Key)
	return run.Config.RemoveLocalConfigValue(step.Key)
}
