package step

type QueueMessage struct {
	Message string
	Empty
}

func (step *QueueMessage) Run(args RunArgs) error {
	args.Runner.FinalMessages.Add(step.Message)
	return nil
}
