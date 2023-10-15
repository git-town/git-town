package opcode

import "github.com/git-town/git-town/v9/src/vm/shared"

type QueueMessage struct {
	Message string
	undeclaredOpcodeMethods
}

func (step *QueueMessage) Run(args shared.RunArgs) error {
	args.Runner.FinalMessages.Add(step.Message)
	return nil
}
