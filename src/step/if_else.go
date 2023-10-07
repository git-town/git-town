package step

type IfElse struct {
	Condition  ConditionFunc
	TrueSteps  []Step // the steps to execute if the given branch is empty
	FalseSteps []Step // the steps to execute if the given branch is not empty
	Empty
}

type ConditionFunc func() (bool, error)

func (step *IfElse) Run(args RunArgs) error {
	condition, err := step.Condition()
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
