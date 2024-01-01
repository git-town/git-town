package opcode

import "github.com/git-town/git-town/v11/src/vm/shared"

type QueueMessage struct {
	Message string
	undeclaredOpcodeMethods
}

func (self *QueueMessage) Run(args shared.RunArgs) error {
	args.Runner.FinalMessages.Add(self.Message)
	return nil
}
