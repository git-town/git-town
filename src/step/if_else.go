package step

import "github.com/git-town/git-town/v9/src/git"

type IfElse struct {
	Condition  ConditionFunc
	TrueSteps  []Step // the steps to execute if the given branch is empty
	FalseSteps []Step // the steps to execute if the given branch is not empty
	Empty
}

type ConditionFunc func(*git.BackendCommands) (bool, error)

func (step *IfElse) Run(args RunArgs) error {
	condition, err := step.Condition(&args.Runner.Backend)
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
