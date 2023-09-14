package steps

import (
	"github.com/git-town/git-town/v9/src/git"
)

type StashOpenChangesStep struct {
	EmptyStep
}

func (step *StashOpenChangesStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&RestoreOpenChangesStep{}}, nil
}

func (step *StashOpenChangesStep) Run(args RunArgs) error {
	return args.Run.Frontend.Stash()
}
