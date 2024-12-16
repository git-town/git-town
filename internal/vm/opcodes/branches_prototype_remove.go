package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/messages"
	"github.com/git-town/git-town/v17/internal/vm/shared"
)

// removes the branch with the given name from the prototype branches list in the Git config
type BranchesPrototypeRemove struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchesPrototypeRemove) Run(args shared.RunArgs) error {
	var err error
	if args.Config.Value.NormalConfig.IsPrototypeBranch(self.Branch) {
		args.FinalMessages.Add(fmt.Sprintf(messages.PrototypeRemoved, self.Branch))
		err = args.Config.Value.NormalConfig.RemoveFromPrototypeBranches(self.Branch)
	}
	return err
}
