package opcodes

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

type SnapshotInitialUpdateLocalSHAIfNeeded struct {
	Branch gitdomain.LocalBranchName
}

func (self *SnapshotInitialUpdateLocalSHAIfNeeded) Run(args shared.RunArgs) error {
	newSHA, err := args.Git.SHAForBranch(args.Backend, self.Branch.BranchName())
	if err != nil {
		return err
	}
	args.PrependOpcodes(&SnapshotInitialUpdateLocalSHA{
		Branch: self.Branch,
		SHA:    newSHA,
	})
	return nil
}
