package opcode

import "github.com/git-town/git-town/v9/src/vm/shared"

type QueueMessage struct {
	Message string
	undeclaredOpcodeMethods
}

func (op *QueueMessage) Run(args shared.RunArgs) error {
	args.Runner.FinalMessages.Add(op.Message)
	return nil
}
