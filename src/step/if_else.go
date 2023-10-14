package step

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/git"
)

// IfElse allows running different steps based on a condition evaluated at runtime.
type IfElse struct {
	Condition  func(*git.BackendCommands, config.Lineage) (bool, error)
	TrueSteps  []Step // the steps to execute if the given branch is empty
	FalseSteps []Step // the steps to execute if the given branch is not empty
	Empty
}

func (step *IfElse) Run(args RunArgs) error {
	condition, err := step.Condition(&args.Runner.Backend, args.Lineage)
	if err != nil {
		return err
	}
	if condition {
		args.AddSteps(step.TrueSteps...)
	} else {
		args.AddSteps(step.FalseSteps...)
	}
	return nil
}
