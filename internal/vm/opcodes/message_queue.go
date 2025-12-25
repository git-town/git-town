package opcodes

import "github.com/git-town/git-town/v22/internal/vm/shared"

type MessageQueue struct {
	Message string
}

func (self *MessageQueue) Run(args shared.RunArgs) error {
	args.FinalMessages.Add(self.Message)
	return nil
}
