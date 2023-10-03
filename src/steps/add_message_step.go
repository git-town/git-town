package steps

type AddMessageStep struct {
	Message string
	EmptyStep
}

func (step *AddMessageStep) Run(args RunArgs) error {
	args.Runner.Stats.RegisterMessage(step.Message)
	return nil
}
