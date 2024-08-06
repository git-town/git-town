package opcodes

import "github.com/git-town/git-town/v15/internal/vm/shared"

type QueueMessage struct {
	Message                 string
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *QueueMessage) Run(args shared.RunArgs) error {
	args.FinalMessages.Add(self.Message)
	return nil
}
