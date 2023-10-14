package opcode

type QueueMessage struct {
	Message string
	undeclaredOpcodeMethods
}

func (step *QueueMessage) Run(args RunArgs) error {
	args.Runner.FinalMessages.Add(step.Message)
	return nil
}
