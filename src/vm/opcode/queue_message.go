package opcode

type QueueMessage struct {
	Message string
	BaseOpcode
}

func (step *QueueMessage) Run(args RunArgs) error {
	args.Runner.FinalMessages.Add(step.Message)
	return nil
}
