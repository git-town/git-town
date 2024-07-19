package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// RemoveFromPrototypeBranches removes the branch with the given name as a prototype branch.
type RemoveFromPrototypeBranches struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RemoveFromPrototypeBranches) Run(args shared.RunArgs) error {
	args.FinalMessages.Add(fmt.Sprintf(messages.PrototypeRemoved, self.Branch))
	return args.Config.RemoveFromPrototypeBranches(self.Branch)
}
