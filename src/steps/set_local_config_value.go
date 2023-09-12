package steps

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

type SetLocalConfigStep struct {
	Key   config.Key
	Value string
	EmptyStep
}

func (step *SetLocalConfigStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&RemoveLocalConfigStep{Key: step.Key}}, nil
}

func (step *SetLocalConfigStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	return run.Config.SetLocalConfigValue(step.Key, step.Value)
}
