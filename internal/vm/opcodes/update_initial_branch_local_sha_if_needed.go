package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

type UpdateInitialBranchLocalSHAIfNeeded struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *UpdateInitialBranchLocalSHAIfNeeded) Run(args shared.RunArgs) error {
	newSHA, err := args.Git.SHAForBranch(args.Backend, self.Branch.BranchName())
	if err != nil {
		return err
	}
	args.PrependOpcodes(&UpdateInitialBranchLocalSHA{
		Branch: self.Branch,
		SHA:    newSHA,
	})
	return nil
}
