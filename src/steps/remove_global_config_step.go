package steps

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

type RemoveGlobalConfigStep struct {
	Key           config.Key // the config key to remove
	previousValue string
	EmptyStep
}

func (step *RemoveGlobalConfigStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&AddGlobalConfigStep{
		Key:   step.Key,
		Value: step.previousValue,
	}}, nil
}

func (step *RemoveGlobalConfigStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	step.previousValue = run.Config.GlobalConfigValue(step.Key)
	return run.Config.RemoveGlobalConfigValue(step.Key)
}
