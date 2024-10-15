package opcodes

import "github.com/git-town/git-town/v16/internal/vm/shared"

type MessageQueue struct {
	Message                 string
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *MessageQueue) Run(args shared.RunArgs) error {
	args.FinalMessages.Add(self.Message)
	return nil
}
