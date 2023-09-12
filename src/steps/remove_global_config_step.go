package steps

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

type RemoveGlobalConfigStep struct {
	Key config.Key // the config key to remove
	EmptyStep
}

func (step *RemoveGlobalConfigStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&AddGlobalConfigStep{}}, nil
}

func (step *RemoveGlobalConfigStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	return run.Config.RemoveGlobalConfigValue(step.Key)
}
