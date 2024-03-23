package opcodes

import (
	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/git-town/git-town/v13/src/vm/shared"
)

type UpdateInitialBranchLocalSHA struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *UpdateInitialBranchLocalSHA) Run(args shared.RunArgs) error {
	newSHA, err := args.Runner.Backend.SHAForBranch(self.Branch.BranchName())
	if err != nil {
		return err
	}
	return args.UpdateInitialBranchLocalSHA(self.Branch, newSHA)
}
