package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

type SnapshotInitialUpdateLocalSHA struct {
	Branch                  gitdomain.LocalBranchName
	SHA                     gitdomain.SHA
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *SnapshotInitialUpdateLocalSHA) Run(args shared.RunArgs) error {
	return args.UpdateInitialSnapshotLocalSHA(self.Branch, self.SHA)
}
