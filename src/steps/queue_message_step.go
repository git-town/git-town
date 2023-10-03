package steps

type QueueMessageStep struct {
	Message string
	EmptyStep
}

func (step *QueueMessageStep) Run(args RunArgs) error {
	args.Runner.Stats.RegisterMessage(step.Message)
	return nil
}
