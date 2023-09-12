package steps

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

type AddGlobalConfigStep struct {
	Key   config.Key
	Value string
	EmptyStep
}

func (step *AddGlobalConfigStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&RemoveGlobalConfigStep{Key: step.Key}}, nil
}

func (step *AddGlobalConfigStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	return run.Config.SetGlobalConfigValue(step.Key, step.Value)
}
