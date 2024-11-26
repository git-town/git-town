package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// BranchLocalDeleteContent deletes the branch with the given name.
type BranchLocalDeleteContent struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchLocalDeleteContent) Run(args shared.RunArgs) error {
	err := args.Git.RebaseOnto(args.Frontend, self.Branch.BranchName(), args.Config.Value.ValidatedConfigData.MainBranch)
	if err != nil {
		return err
	}
	return args.Git.DeleteLocalBranch(args.Frontend, self.Branch)
}
