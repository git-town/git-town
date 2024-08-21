package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/messages"
	"github.com/git-town/git-town/v15/internal/vm/shared"
)

// removes the branch with the given name from the prototype branches list in the Git config
type RemoveFromPrototypeBranches struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RemoveFromPrototypeBranches) Run(args shared.RunArgs) error {
	args.FinalMessages.Add(fmt.Sprintf(messages.PrototypeRemoved, self.Branch))
	return args.Config.RemoveFromPrototypeBranches(self.Branch)
}
